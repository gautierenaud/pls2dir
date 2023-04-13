package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gautierenaud/pls2dir/internal/parsers"
)

var (
	PlaylistPath      string
	selectedPlaylists []string
	destination       string

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
	exportCmd.Flags().StringSliceVarP(&selectedPlaylists, "selection", "s", nil, "playlist to export (comma separated)")
	exportCmd.Flags().StringVarP(&destination, "destination", "d", "", "destination folder to export the playlist to")
	exportCmd.MarkFlagRequired("destination")

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

	// TODO: use a io.TeeReader to put source's content in RAM only once
	for _, file := range files {
		filename := path.Base(file)
		destinationPath := path.Join(destination, filename)

		export, err := ShouldExport(file, destinationPath)
		if err != nil {
			return err
		}

		if !export {
			log.Println("Skipping", filename)
			continue
		}

		log.Println("Exporting", file)

		source, err := os.Open(file)
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
	}

	return nil
}

func ShouldExport(source, dest string) (bool, error) {
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

	sourceFile, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer sourceFile.Close()

	sourceHash := sha256.New()
	_, err = io.Copy(sourceHash, sourceFile)
	if err != nil {
		return false, err
	}

	destFile, err := os.Open(dest)
	if err != nil {
		return false, err
	}
	defer destFile.Close()

	destHash := sha256.New()
	_, err = io.Copy(destHash, destFile)
	if err != nil {
		return false, err
	}

	if hex.EncodeToString(sourceHash.Sum(nil)) == hex.EncodeToString(destHash.Sum(nil)) {
		return false, nil
	}

	return true, nil
}
