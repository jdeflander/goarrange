package index_test

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jdeflander/goarrange/internal/arranger/internal/index"
)

func TestIndex(t *testing.T) {
	if err := filepath.Walk("testdata", walk); err != nil {
		t.Error(err)
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("failed walking to file at '%s': %w", path, err)
	}
	if !info.IsDir() {
		return nil
	}

	set := token.NewFileSet()
	packages, err := parser.ParseDir(set, path, nil, 0)
	if err != nil {
		return fmt.Errorf("failed parsing directory at '%s': %v", path, err)
	}

	for _, pkg := range packages {
		idx := index.New(pkg)
		for name, file := range pkg.Files {
			gotSort := idx.Sort(file.Decls)
			golden := fmt.Sprintf("%slden.json", name)
			bytes, err := ioutil.ReadFile(golden)
			if err != nil {
				return fmt.Errorf("failed reading golden file at '%s': %v", golden, err)
			}
			var wantSort []int
			if err := json.Unmarshal(bytes, &wantSort); err != nil {
				return fmt.Errorf("failed unmarshalling golden file at '%s': %v", golden, err)
			}
			if diff := cmp.Diff(gotSort, wantSort); diff != "" {
				return fmt.Errorf("sort mismatch for file at '%s' (-got +want):\n%s", name, diff)
			}

			gotSorted := idx.Sorted(file.Decls)
			wantSorted := sort.IntsAreSorted(wantSort)
			if diff := cmp.Diff(gotSorted, wantSorted); diff != "" {
				return fmt.Errorf("sorted mismatch for file at '%s' (-got +want):\n%s", name, diff)
			}
		}
	}
	return nil
}
