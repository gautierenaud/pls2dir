package parsers

import (
	"net/url"
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

	res := make(map[string][]string)
	name := strings.TrimSuffix(filepath.Base(playlistPath), filepath.Ext(playlistPath))
	for _, key := range cfg.Section("playlist").KeyStrings() {
		if strings.HasPrefix(key, "File") {
			val := cfg.Section("playlist").Key(key)
			if val == nil {
				// should never happen
				continue
			}
			u, _ := url.ParseRequestURI(val.MustString(""))
			res[name] = append(res[name], u.Path)
		}
	}

	return res, nil
}
