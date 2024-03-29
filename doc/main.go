package main

import (
	"log"

	"github.com/asymmetricia/ghastly/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.Root, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}
