package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "scan" {
		if flag.Arg(1) == "" {
			fmt.Printf("%s", "You must provide directory to scan, exiting...\n")
			os.Exit(1)
		}
		db := dbInit()
		dir := flag.Arg(1)
		scan(dir, visit(db))
		os.Exit(0)
	} else if flag.Arg(0) == "serve" {
		db := dbInit()
		server(db)
	} else {
		fmt.Printf("%s", "You must use scan or serve command, exiting...\n")
		os.Exit(1)
	}
}
