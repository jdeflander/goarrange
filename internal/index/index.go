package index

import (
	"go/ast"
	"sort"
)

type Index struct {
	decls map[ast.Decl]int
}

func New() Index {
	decls := map[ast.Decl]int{}
	return Index{decls: decls}
}

func (i Index) Append(decl ast.Decl) {
	i.decls[decl] = len(i.decls)
}

func (i Index) IsSorted(decls []ast.Decl) bool {
	records := i.records(decls)
	return sort.IsSorted(records)
}

func (i Index) Sort(decls []ast.Decl) []int {
	records := i.records(decls)
	sort.Stable(records)

	var indexes []int
	for _, record := range records {
		indexes = append(indexes, record.index)
	}
	return indexes
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
