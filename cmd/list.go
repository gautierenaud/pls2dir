package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/gautierenaud/pls2dir/internal/parsers"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists playlists that can be exported",
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	p := parsers.NewRhythmDBParser()

	playlists, err := p.ParsePlaylist(PlaylistPath)
	if err != nil {
		return err
	}

	log.Printf("Available playlists for file %s are:\n", PlaylistPath)
	for name := range playlists {
		fmt.Println(name)
	}

	return nil
}
