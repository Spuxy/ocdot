package main

import (
	"flag"
)


func main() {
	var pathToDots string
	flag.StringVar(
		&pathToDots,
		"path to your dots",
		"./dotfiles",
		"",
	)
	flag.Parse()
}
