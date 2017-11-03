# GoTag

Privacy focused, browser based file manager for tagging and searching local files on disk.

Currently this project is only prototype and have many issues to fix, don't use it in production please.

You may safely test it even with important files because it don't make any changes to files it handles.
Only changes are made in app's folder.

## How to use

### Ubuntu 16.04

- install packages by using command:
```bash
# packages required to make .webm and .mp4 thumbnails
apt-get install -y libavutil-dev libavformat-dev libswscale-dev
```
- download ready to use archive from [https://github.com/SystemZ/gotag/releases](https://github.com/SystemZ/gotag/releases)
- unpack archive
- run binary to scan dir add add files to gotag.sqlite3 DB located in same directory as app
```
./gotag scan <your dir with files>
```
- run app for web interface at http://127.0.0.1:3000
```
./gotag serve
```
- try to enjoy besides many bugs

## Known issues

Currently thumbnails are created in-the-fly.
Some part of thumbnail handle code makes memory leak by using `/img/thumb` url.

Probably something is wrong with release of resources after `.Decode`
```go
jpeg.Decode(imgFile)
```

Easy fix is to create thumbnails on scan.

## TODO

Create thumbnails on scan, located in app folder with folder hierarchy based on sha256 of file.

Example:

```
- gotag (dir where app starts)
  - cache
    - 1c
      - 23
        - 55
          - 7e1514d90f8db310ec55de98315d267f5c7cadcf021a2507415498fc2b
```

## License

MIT