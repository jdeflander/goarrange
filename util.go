package main

import (
	"fmt"
	"io"
	"os"
)

func failWithUsage() {
	writeUsage(os.Stderr)
	os.Exit(1)
}

func fprintf(writer io.Writer, format string, args ...interface{}) {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		panic(err)
	}
}

func writeUsage(writer io.Writer) {
	usage := `Automatic arrangement of Go source code

Usage:
  %[1]s help
  %[1]s run
`
	name := os.Args[0]
	fprintf(writer, usage, name)
}
