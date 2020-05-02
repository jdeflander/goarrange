package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jdeflander/goarrange/internal/arranger"
)

func arrange(directory, filename string, dryRun bool) error {
	set := token.NewFileSet()
	packages, err := parser.ParseDir(set, directory, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed parsing directory at '%s': %v", directory, err)
	}

	for _, pkg := range packages {
		a := arranger.New(pkg, set)
		for path, file := range pkg.Files {
			if filename != "" && filepath.Base(path) != filename || a.Arranged(file) {
				continue
			}

			if dryRun {
				if _, err := fmt.Println(path); err != nil {
					return fmt.Errorf("failed printing path '%s': %w", path, err)
				}
			} else {
				src, err := ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed reading file at '%s': %v", path, err)
				}
				dst := a.Arrange(file, src)
				if err := ioutil.WriteFile(path, dst, 0644); err != nil {
					return fmt.Errorf("failed writing file  at '%s': %v", path, err)
				}
			}
		}
	}
	return nil
}

func split(path string) (string, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", fmt.Errorf("failed checking status of file at '%s': %w", path, err)
	}
	if info.IsDir() {
		return path, "", nil
	}
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	return dir, filename, nil
}
