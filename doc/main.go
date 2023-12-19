package main

import (
	"log"

	"github.com/pdbogen/ghastly/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.Root, "./")
	if err != nil {
		log.Fatal(err)
	}
}
