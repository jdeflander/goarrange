package arranger_test

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jdeflander/goarrange/internal/arranger"
)

func TestArranger(t *testing.T) {
	if err := filepath.Walk("testdata", walk); err != nil {
		t.Error(err)
	}
}

func notGolden(info os.FileInfo) bool {
	name := info.Name()
	return !strings.HasSuffix(name, ".golden.go")
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("failed walking to file at '%s': %w", path, err)
	}
	if !info.IsDir() {
		return nil
	}

	set := token.NewFileSet()
	packages, err := parser.ParseDir(set, path, notGolden, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed parsing directory at '%s': %v", path, err)
	}

	for _, pkg := range packages {
		a := arranger.New(pkg, set)
		for name, file := range pkg.Files {
			bs, err := ioutil.ReadFile(name)
			if err != nil {
				return fmt.Errorf("failed reading file at '%s': %v", name, err)
			}
			gotArrange := a.Arrange(file, bs)
			golden := fmt.Sprintf("%slden.go", name)
			wantArrange, err := ioutil.ReadFile(golden)
			if err != nil {
				return fmt.Errorf("failed reading golden file at '%s': %v", golden, err)
			}
			if diff := cmp.Diff(gotArrange, wantArrange); diff != "" {
				return fmt.Errorf("arrange output mismatch for file at '%s':\n%s", name, diff)
			}

			gotArranged := a.Arranged(file)
			wantArranged := bytes.Equal(bs, gotArrange)
			if diff := cmp.Diff(gotArranged, wantArranged); diff != "" {
				return fmt.Errorf("arranged output mismatch for file at '%s':\n%s", name, diff)
			}
		}
	}
	return nil
}
