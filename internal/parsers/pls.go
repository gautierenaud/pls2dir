package parsers

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type PlsParser struct{}

func NewPlsParser() PlaylistParser {
	return PlsParser{}
}

func (p PlsParser) ParsePlaylist(playlistPath string) (map[string][]string, error) {
	cfg, err := ini.Load(playlistPath)
	if err != nil {
		return nil, err
	}

	playlistDir := path.Dir(playlistPath)

	res := make(map[string][]string)
	name := strings.TrimSuffix(filepath.Base(playlistPath), filepath.Ext(playlistPath))
	for _, key := range cfg.Section("playlist").KeyStrings() {
		if strings.HasPrefix(key, "File") {
			val := cfg.Section("playlist").Key(key)
			if val == nil {
				// should never happen
				continue
			}

			rawPath := val.MustString("")

			var filepath string
			// when the file is in another drive, it might start with "file://"
			u, err := url.ParseRequestURI(rawPath)
			if err != nil {
				filepath = path.Join(playlistDir, rawPath)
			} else {
				filepath = u.Path
			}

			res[name] = append(res[name], filepath)
		}
	}

	return res, nil
}
