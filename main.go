package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/godoes/winseq" // Windows 虚拟终端序列
	"github.com/jessevdk/go-flags"
)

// CLI utility for converting markdown to a single html file
func main() {
	var opts Options
	inputs, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if len(opts.InputFile) > 0 {
		inputs = []string{opts.InputFile}
	}
	if len(inputs) <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Please specify input Markdown")
		os.Exit(1)
	}

	var files []string
	for _, input := range inputs {
		f, err := filepath.Glob(input)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		files = append(files, f...)
	}
	if len(files) <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "File is not found")
		os.Exit(1)
	}

	parser := HTMLParser{Options: opts}
	parser.parserMarkdown(files)
}
