package core

import (
	"log"
	"os"

	"github.com/carlogit/phash"
)

// https://stackoverflow.com/a/38469006/1351857
func GetPHash(filename string) string {
	img, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	hash, err := phash.GetHash(img)
	if err != nil {
		log.Fatal(err)
	}
	return hash
}
