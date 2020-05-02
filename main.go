package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const usage = `Automatic arrangement of Go source code

Usage:
  %[1]s help
  %[1]s run [-d] [-p=<path>] [-r]

Options:
  -d        Dry-run listing unarranged files
  -p=<path> Path of file or directory to arrange [default: .]
  -r        Walk directories recursively
`

func main() {
	log.SetFlags(0)
	path := os.Args[0]
	name := filepath.Base(path)
	prefix := fmt.Sprintf("%s: ", name)
	log.SetPrefix(prefix)

	if len(os.Args) < 2 {
		log.Fatalf("missing command, run '%s help' for usage information", name)
	}

	switch os.Args[1] {
	case "help":
		if _, err := fmt.Printf(usage, name); err != nil {
			log.Fatalf("failed printing usage information: %v", err)
		}

	case "run":
		set := flag.NewFlagSet("", flag.ContinueOnError)
		var dryRun bool
		set.BoolVar(&dryRun, "d", false, "")
		var path string
		set.StringVar(&path, "p", ".", "")
		var recursive bool
		set.BoolVar(&recursive, "r", false, "")
		set.SetOutput(ioutil.Discard)

		args := os.Args[2:]
		if err := set.Parse(args); err != nil {
			log.Fatalf("invalid options, run '%s help' for usage information", name)
		}

		directory, filename, err := split(path)
		if err != nil {
			log.Fatalf("failed splitting path: %v", err)
		}

		if filename == "" && recursive {
			walk := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("failed walking to file at '%s': %w", path, err)
				}
				if !info.IsDir() {
					return nil
				}

				if err := arrange(path, "", dryRun); err != nil {
					return fmt.Errorf("failed arranging directory at '%s': %w", path, err)
				}
				return nil
			}
			if err := filepath.Walk(directory, walk); err != nil {
				log.Fatalf("failed walking directory at '%s': %v", directory, err)
			}
		} else if err := arrange(directory, filename, dryRun); err != nil {
			log.Fatalf("failed arranging directory at '%s': %v", directory, err)
		}

	default:
		log.Fatalf("invalid command, run '%s help' for usage information", name)
	}
}
