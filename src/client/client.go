package main

import (
	"log"
	"os"
	"sdfs/client/client_util"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("issue")
	}

	command := os.Args[1]
	if len(os.Args) == 2 && command == "ls" {
		client_util.Ls()
		os.Exit(0)
	}

	file := os.Args[2]
	if command == "create" {
		client_util.Create(file)
	} else if command == "read" {
		client_util.Read(file)
	} else if command == "rm" {
		client_util.Rm(file)
	}
	os.Exit(0)
}
