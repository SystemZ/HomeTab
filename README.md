# GoTag

Privacy focused, browser based file manager for tagging and searching local files on disk.

Currently this project is only prototype and have many issues to fix, don't use it in production please.

You may safely test it even with important files because it don't make any changes to files it handles.
Only changes are made in app's folder.

## How to use

### Ubuntu 16.04

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

## License

MIT