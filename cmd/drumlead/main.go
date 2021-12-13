package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/peter-mueller/drumlead/parser"
	"github.com/peter-mueller/drumlead/render"
)

func main() {
	for _, file := range os.Args[1:] {
		extension := filepath.Ext(file)
		if extension != ".md" {
			log.Fatal(fmt.Errorf("only .md files are allowed, but was %s (%s)", file, extension))
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to open file %s: %w", file, err))
		}
		l := parser.Parse(string(data))

		var (
			name    = file[:len(extension)+1]
			pdfName = fmt.Sprintf("%s.pdf", name)
		)
		render.SaveFile(l, pdfName)
	}

}
