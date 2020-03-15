package main

import (
	"flag"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		failWithUsage()
	}

	switch os.Args[1] {
	case "help":
		help()

	case "run":
		set := flag.NewFlagSet("", flag.ContinueOnError)
		path := set.String("p", ".", "")
		set.SetOutput(ioutil.Discard)

		args := os.Args[2:]
		if err := set.Parse(args); err != nil {
			failWithUsage()
		}

		if err := run(*path); err != nil {
			fprintf(os.Stderr, "failed running: %s\n", err)
			os.Exit(1)
		}

	default:
		failWithUsage()
	}
}
