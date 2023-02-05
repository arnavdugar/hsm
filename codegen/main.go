package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/arnavdugar/hsm/codegen/golang"
	"github.com/arnavdugar/hsm/codegen/mermaid"
	"github.com/arnavdugar/hsm/codegen/parser"
)

var inFileName = flag.String("i", "", "Input file.")
var outFileName = flag.String("o", "", "Output file.")

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	if inFileName == nil || *inFileName == "" {
		return errors.New("input file is required")
	}

	if outFileName == nil || *outFileName == "" {
		return errors.New("output file is required")
	}

	var inFile *os.File
	if *inFileName == "-" {
		inFile = os.Stdin
	} else {
		var err error
		inFile, err = os.Open(*inFileName)
		if err != nil {
			return err
		}
	}

	machine, err := parser.Parse(inFile)
	if err != nil {
		return err
	}

	golangGenerator := golang.Create(machine)
	output, err := golangGenerator.Render()
	if err != nil {
		return err
	}

	var outFile *os.File
	if *outFileName == "-" {
		outFile = os.Stdout
	} else {
		var err error
		outFile, err = os.Create(*outFileName)
		if err != nil {
			return err
		}
	}

	_, err = outFile.Write(output)
	if err != nil {
		return err
	}

	if machine.Machine.Codegen.Mermaid.Enabled {
		mermaidGenerator := mermaid.Create(machine)
		output, err = mermaidGenerator.Render()
		if err != nil {
			return err
		}

		outFile, err = os.Create(machine.Machine.Codegen.Mermaid.Filename)
		if err != nil {
			return err
		}

		_, err = outFile.Write(output)
		if err != nil {
			return err
		}
	}

	return nil
}
