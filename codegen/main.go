package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arnavdugar/hsm/codegen/golang"
	"github.com/arnavdugar/hsm/codegen/mermaid"
	"github.com/arnavdugar/hsm/codegen/parser"
	"github.com/arnavdugar/hsm/codegen/typescript"
)

type output struct {
	target string
	path   string
}

type outputs []output

func (items *outputs) add(spec string) error {
	target, path, ok := strings.Cut(spec, "=")
	if !ok || path == "" {
		return errors.New("expected target=path")
	}
	switch target {
	case "go", "typescript", "mermaid":
	default:
		return fmt.Errorf("unknown output target %q", target)
	}
	for _, item := range *items {
		if item.target == target {
			return fmt.Errorf("output target %q is specified more than once", target)
		}
		if filepath.Clean(item.path) == filepath.Clean(path) {
			return fmt.Errorf("output path %q is specified more than once", path)
		}
	}
	*items = append(*items, output{target: target, path: path})
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil && !errors.Is(err, flag.ErrHelp) {
		log.Fatal(err)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("hsm", flag.ContinueOnError)
	input := flags.String("i", "", "Input YAML file, or - for standard input.")
	var targets outputs
	flags.Func("o", "Output target and path (go=path, typescript=path, mermaid=path). Repeat to emit multiple targets. Use - for standard output.", targets.add)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if *input == "" {
		return errors.New("input file is required")
	}

	if len(targets) == 0 {
		return errors.New("specify at least one output with -o target=path")
	}

	inFile := os.Stdin
	if *input != "-" {
		file, err := os.Open(*input)
		if err != nil {
			return err
		}
		defer file.Close()
		inFile = file
	}
	machine, err := parser.Parse(inFile)
	if err != nil {
		return err
	}

	// Render every target before opening any output file.
	rendered := make([][]byte, len(targets))
	for i, item := range targets {
		switch item.target {
		case "go":
			rendered[i], err = golang.Create(machine).Render()
		case "typescript":
			rendered[i], err = typescript.Create(machine).Render()
		case "mermaid":
			rendered[i], err = mermaid.Create(machine).Render()
		}
		if err != nil {
			return fmt.Errorf("render %s: %w", item.target, err)
		}
	}
	for i, item := range targets {
		if item.path == "-" {
			_, err = os.Stdout.Write(rendered[i])
		} else {
			err = os.WriteFile(item.path, rendered[i], 0666)
		}
		if err != nil {
			return fmt.Errorf("write %s: %w", item.target, err)
		}
	}
	return nil
}
