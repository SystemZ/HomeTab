# GoTag

Privacy focused, browser based file manager for tagging and searching local files on disk.

You may safely test it even with important files because it don't make any changes to files it handles.
Only changes are made in MySQL DB and thumbnails cache dir

## How to use

### Linux

- run cmd to scan dir and add files to DB:
```
./gotag scan <your dir with files>
```
- run app for web interface at http://127.0.0.1:3000
```
./gotag serve
```

You can also use built-in scan from WebGUI but it's in raw state, without progress bar.  
If you use it, please observe stdout logs for progress info

### Windows

This app is primary for Linux powered NAS devices.  
It may run on Windows but it's not supported by author

### Platform differences

Currently webm/mp4/gif thumbnails are supported only on Linux amd64 build as a limitation of [lilliput](https://github.com/discordapp/lilliput) library

## Notes

### Deploy notes

Try to use best practices for image:
* https://github.com/chemidy/smallest-secured-golang-docker-image/blob/master/Dockerfile

### Dev notes

* https://github.com/discordapp/lilliput/issues/55
* https://github.com/discordapp/lilliput/issues/24

Lilliput lib have problems with build when go mod is used.  

If you are getting something like this:
```plain
  go get github.com/discordapp/lilliput
go: finding github.com/discordapp/lilliput latest
go: extracting github.com/discordapp/lilliput v0.0.0-20191204003513-dd93dff726a5
# github.com/discordapp/lilliput
/usr/bin/ld: cannot find -lpng
/usr/bin/ld: cannot find -lpng
collect2: error: ld returned 1 exit status

```

you need to enter dir like this:
```plain
$GOPATH/pkg/mod/github.com/discordapp/lilliput@v0.0.0-20191204003513-dd93dff726a5/deps/linux/lib
```

and manually copy two files, by default links doesn't work for this
```bash
sudo cp libpng16.a libpng.a
sudo cp libpng16.la libpng.la
```

Next run with `go build` or `go get` should be ok
