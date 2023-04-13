# pls2dir
Exports playlists to directories

TODOs:
* recreate with some folders dept (like artist_name/music.file)
    * enable export musics with the same name
* progression bar
* skip missing file flag
* skip/error counter or list
* io.TeeReader when reading from source

# Supported files

## .pls

INI like format, as described [here](https://en.wikipedia.org/wiki/PLS_(file_format)).

## .xml

Rhythmbox's internal representation of playlists.