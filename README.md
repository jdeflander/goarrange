# goarrange

Ever wanted a consistent ordering for declarations in your Go code? With `goarrange`, you can automatically follow the
conventions of GoDoc! Constants come first, followed by variables, functions and types with their associated constants,
variables, functions and methods. Within each of these categories, exported declarations precede unexported ones. Lastly
`goarrange` enforces an alphabetic ordering.

## Installation

```sh
$ go get github.com/jdeflander/goarrange
```

### v1.17 and later

```sh
$ go install github.com/jdeflander/goarrange@v1.0.0
```

## Usage

```sh
$ goarrange help
Automatic arrangement of Go source code

Usage:
  goarrange help
  goarrange run [-d] [-p=<path>] [-r]

Options:
  -d        Dry-run listing unarranged files
  -p=<path> Path of file or directory to arrange [default: .]
  -r        Walk directories recursively
```
