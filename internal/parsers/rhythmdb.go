package parsers

import (
	"encoding/xml"
	"io"
	"net/url"
	"os"
)

type RhythmDBPlaylists struct {
	Playlists []RhythmDBPlaylist `xml:"playlist"`
}

type RhythmDBPlaylist struct {
	Name      string   `xml:"name,attr"`
	Locations []string `xml:"location"`
}

type RhythmDBParser struct{}

func NewRhythmDBParser() PlaylistParser {
	return RhythmDBParser{}
}

func (p RhythmDBParser) ParsePlaylist(playlistPath string) (map[string][]string, error) {
	playlistFile, err := os.Open(playlistPath)
	if err != nil {
		return nil, err
	}
	defer playlistFile.Close()

	playlistRaw, err := io.ReadAll(playlistFile)
	if err != nil {
		return nil, err
	}

	var playlists RhythmDBPlaylists
	err = xml.Unmarshal(playlistRaw, &playlists)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	for _, p := range playlists.Playlists {
		if len(p.Locations) > 0 {
			res[p.Name] = make([]string, 0, len(p.Locations))

			for _, location := range p.Locations {
				u, _ := url.ParseRequestURI(location)
				res[p.Name] = append(res[p.Name], u.Path)
			}
		}
	}

	return res, nil
}
