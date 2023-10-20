package main

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
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
		_, _ = colorRed.Fprintln(color.Error, "Please specify input Markdown")
		os.Exit(1)
	}

	var files []string
	for _, input := range inputs {
		var f []string
		if f, err = filepath.Glob(input); err != nil {
			_, _ = colorRed.Fprintln(color.Error, err)
			os.Exit(1)
		}
		files = append(files, f...)
	}
	if len(files) <= 0 {
		_, _ = colorRed.Fprintln(color.Error, "File is not found")
		os.Exit(1)
	}

	parser := HTMLParser{Options: opts}
	parser.parserMarkdown(files)
}
