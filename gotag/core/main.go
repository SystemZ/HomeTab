package core

import (
	"gitlab.com/systemz/gotag/config"
	"strconv"
)

func thumbDirPath(sha256 string) (path string) {
	parent := "./cache"
	if config.CACHE_DIR != "" {
		parent = config.CACHE_DIR
	}
	lvl1 := string(sha256[0]) + string(sha256[1])
	lvl2 := string(sha256[2]) + string(sha256[3])
	lvl3 := string(sha256[4]) + string(sha256[5])
	path = parent + "/" + lvl1 + "/" + lvl2 + "/" + lvl3 + "/"
	return path
}

func ThumbPath(sha256 string, width uint, height uint) (path string) {
	dir := thumbDirPath(sha256)
	path = dir + sha256 + "_" + strconv.Itoa(int(width)) + "_" + strconv.Itoa(int(height))
	return path
}
