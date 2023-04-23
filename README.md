# pls2dir

Exports playlists to directories

Usage:
```
pls2dir -p <path_to_playlist> -d <destination> --folder-cutoff <prefix_to_trim>
```

`--folder-cutoff` is used when you want to reproduce a specific folder structure at your destination.
E.g. you have organized your music library in the following way:
```
<path_to_your_library>
├── Artist1
│   └── music.mp3
└── Artist2
    └── AlbumA
        ├── music.flac
        └── music.wav
```
The full path of your music file would become `<path_to_your_library>/Artist2/AlbumA/music.flac`. Now if you want to reproduce a similar structure you can just pass `--folder-cutoff <path_to_your_library>` and there you go.

# Supported Formats

## .pls

INI like format, as described [here](https://en.wikipedia.org/wiki/PLS_(file_format)).

## .xml

Rhythmbox's internal representation of playlists.