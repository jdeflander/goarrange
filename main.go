package main

import "os"

func main() {
	if len(os.Args) < 2 {
		failWithUsage()
	}

	switch os.Args[1] {
	case "help":
		help()

	case "run":
		if err := run(); err != nil {
			fprintf(os.Stderr, "failed running: %s\n", err)
			os.Exit(1)
		}

	default:
		failWithUsage()
	}
}
