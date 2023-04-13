package parsers

import "path/filepath"

type PlaylistParser interface {
	// ParsePlaylist parses a playlist and return a map of path of music files.
	ParsePlaylist(playlist string) (map[string][]string, error)
}

func GetParser(playlist string) PlaylistParser {
	switch filepath.Ext(playlist) {
	case ".xml":
		return NewRhythmboxParser()
	case ".pls":
		return NewPlsParser()
	default:
		return NewPlsParser()
	}
}
