package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	flag.Parse()

	// DB stuff
	db := dbInit()
	allFiles := dbCountAllFiles(db)
	log.Printf("All files in DB: %d \n", allFiles)

	if flag.Arg(0) == "scan" {
		if flag.Arg(1) == "" {
			fmt.Printf("%s", "You must provide directory to scan, exiting...\n")
			os.Exit(1)
		}
		dir := flag.Arg(1)
		generateThumbs, _ := strconv.ParseBool(flag.Arg(2))
		scan(dir, visit(db, generateThumbs))
		os.Exit(0)
	} else if flag.Arg(0) == "serve" {
		server(db)
	} else {
		fmt.Printf("%s", "You must use scan or serve command, exiting...\n")
		os.Exit(1)
	}
}
