package index

import (
	"go/ast"
	"go/doc"
	"sort"
)

// Index represents an index of a Go package's arrangeable declarations.
type Index struct {
	decls map[ast.Decl]int
}

// New returns an index for the given Go package.
func New(pkg *ast.Package) Index {
	decls := map[ast.Decl]int{}
	idx := Index{decls: decls}
	p := doc.New(pkg, "", doc.AllDecls|doc.PreserveAST)

	idx.appendValues(p.Consts)
	idx.appendValues(p.Vars)
	idx.appendFuncs(p.Funcs)

	for _, typ := range p.Types {
		idx.append(typ.Decl)
		idx.appendValues(typ.Consts)
		idx.appendValues(typ.Vars)
		idx.appendFuncs(typ.Funcs)
		idx.appendFuncs(typ.Methods)
	}
	return idx
}

// Sort returns the indices of the given declarations, sorted according to the given index.
func (i Index) Sort(decls []ast.Decl) []int {
	records := i.records(decls)
	sort.Stable(records)

	var indexes []int
	for _, record := range records {
		indexes = append(indexes, record.index)
	}
	return indexes
}

// Sorted checks whether the given declarations are sorted according to the given index.
func (i Index) Sorted(decls []ast.Decl) bool {
	records := i.records(decls)
	return sort.IsSorted(records)
}

func (i Index) append(decl ast.Decl) {
	i.decls[decl] = len(i.decls)
}

func (i Index) appendFuncs(funcs []*doc.Func) {
	for _, fun := range funcs {
		i.append(fun.Decl)
	}
}

func (i Index) appendValues(values []*doc.Value) {
	for _, value := range values {
		i.append(value.Decl)
	}
}

func (i Index) records(decls []ast.Decl) records {
	var records records
	for index, decl := range decls {
		key, ok := i.decls[decl]
		record := record{
			index: index,
			key:   key,
			ok:    ok,
		}
		records = append(records, record)
	}
	return records
}
