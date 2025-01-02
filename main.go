// package ocdot
package main

import (
	"log"
	"os"
)

var pathToDotfiles string

func main() {
	logFile, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)

	rootCmd.PersistentFlags().StringVar(&pathToDotfiles, "path-to-file", ".dotfiles", "path to dotfile (default is $HOME/.dotfiles)")
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error occured: %s", err.Error())
		os.Exit(1)
	}
}
