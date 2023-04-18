package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/gautierenaud/pls2dir/internal/parsers"
)

var (
	PlaylistPath      string
	selectedPlaylists []string
	destination       string
	folderCutoff      string

	exportCmd = &cobra.Command{
		Use:           "pls2dir",
		Short:         "Exports music files listed in a playlist to a directory",
		RunE:          exportPlaylist,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func init() {
	exportCmd.PersistentFlags().StringVarP(&PlaylistPath, "playlist", "p", "", "playlist file to export")
	exportCmd.MarkPersistentFlagRequired("playlist")
	exportCmd.Flags().StringVarP(&destination, "destination", "d", "", "destination folder to export the playlist to")
	exportCmd.MarkFlagRequired("destination")
	exportCmd.Flags().StringSliceVarP(&selectedPlaylists, "selection", "s", nil, "playlist to export (comma separated)")
	exportCmd.Flags().StringVar(&folderCutoff, "folder-cutoff", "", "path prefix to trim when keeping the folder structure")

	exportCmd.AddCommand(ListCmd)
}

func Execute() {
	err := exportCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func exportPlaylist(cmd *cobra.Command, args []string) error {
	p := parsers.GetParser(PlaylistPath)

	playlists, err := p.ParsePlaylist(PlaylistPath)
	if err != nil {
		return err
	}

	totSize := 0
	for _, files := range playlists {
		totSize += len(files)
	}

	files := make([]string, 0, totSize)

	if len(selectedPlaylists) > 0 {
		unknowns := make([]string, 0)
		for _, selected := range selectedPlaylists {
			f, ok := playlists[selected]
			if !ok {
				unknowns = append(unknowns, selected)
			} else {
				files = append(files, f...)
			}
		}

		if len(unknowns) > 0 {
			return fmt.Errorf("got unknown playlist(s) in the selection: %s", strings.Join(unknowns, ", "))
		}
	} else {
		for _, f := range playlists {
			files = append(files, f...)
		}
	}

	log.Printf("Export %d files to %s\n", len(files), destination)

	bar := progressbar.Default(int64(len(files)))
	var skippedFiles []string
	type erroredFile struct {
		file string
		err  error
	}
	var erroredFiles []erroredFile

	// TODO: use a io.TeeReader to put source's content in RAM only once
	for _, file := range files {
		bar.Add(1)

		dir, filename := path.Split(file)

		destinationFolder := destination
		if folderCutoff != "" {
			destinationFolder = path.Join(destination, strings.TrimPrefix(dir, folderCutoff))
		}

		destinationPath := path.Join(destinationFolder, filename)

		export, err := shouldExport(file, destinationPath)
		if err != nil {
			erroredFiles = append(erroredFiles, erroredFile{filename, err})
			continue
		}

		if !export {
			skippedFiles = append(skippedFiles, filename)
			continue
		}

		err = os.MkdirAll(destinationFolder, os.ModePerm)
		if err != nil {
			log.Printf("aaaaaaaa: %s\n", err)
			erroredFiles = append(erroredFiles, erroredFile{filename, err})
			// continue
		}

		source, err := os.Open(file)
		if err != nil {
			erroredFiles = append(erroredFiles, erroredFile{filename, err})
			continue
		}
		defer source.Close()

		destination, err := os.Create(destinationPath)
		if err != nil {
			erroredFiles = append(erroredFiles, erroredFile{filename, err})
			continue
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			erroredFiles = append(erroredFiles, erroredFile{filename, err})
			continue
		}
	}

	if len(skippedFiles) > 0 {
		log.Printf("%d files skipped\n", len(skippedFiles))
	}

	if len(erroredFiles) > 0 {
		log.Printf("%d files skipped due to error:\n", len(erroredFiles))
		for _, erroredFile := range erroredFiles {
			log.Printf("* %s: %s\n", erroredFile.file, erroredFile.err)
		}
	}

	return nil
}

func shouldExport(source, dest string) (bool, error) {
	sourceStat, err := os.Stat(source)
	if err != nil {
		return false, err
	}

	destStat, err := os.Stat(dest)
	if err != nil {
		// missing dest, we want to export
		return true, nil
	}

	if sourceStat.Size() != destStat.Size() {
		// not the same size, so we should probably export
		return true, nil
	}

	// checking for the sha256 took too much time

	return false, nil
}
