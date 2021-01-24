package gotagcore

import (
	"os"

	"github.com/carlogit/phash"
)

// https://stackoverflow.com/a/38469006/1351857
func GetPHash(filename string) (string, error) {
	img, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer img.Close()

	hash, err := phash.GetHash(img)
	if err != nil {
		return hash, err
	}
	return hash, nil
}
