# GoTag
[![Build Status](https://travis-ci.org/SystemZ/gotag.svg?branch=master)](https://travis-ci.org/SystemZ/gotag)

Privacy focused, browser based file manager for tagging and searching local files on disk.

You may safely test it even with important files because it don't make any changes to files it handles.
Only changes are made in app's folder to sqlite DB and thumbnails cache

## How to use

### Linux/Windows

- download ready to use archive from [here](https://github.com/SystemZ/gotag/releases)
- unpack archive
- run binary to scan dir add add files to gotag.sqlite3 DB located in same directory as app
```
./gotag scan <your dir with files>
```
- run app for web interface at http://127.0.0.1:3000
```
./gotag serve
```
- please report encountered issues [here](https://github.com/SystemZ/gotag/issues)

### Platform differences

Currently webm/mp4/gif thumbnails are supported only on Linux amd64 build as a limitation of [lilliput](https://github.com/discordapp/lilliput) library

## License

MIT
