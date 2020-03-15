package main

import (
	"bytes"
	"fmt"
	"github.com/jdeflander/goarrange/internal/index"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

func appendFuncs(idx index.Index, funcs []*doc.Func) {
	for _, fun := range funcs {
		idx.Append(fun.Decl)
	}
}

func appendValues(idx index.Index, values []*doc.Value) {
	for _, value := range values {
		idx.Append(value.Decl)
	}
}

func arrangeFile(file *ast.File, idx index.Index, path string, set *token.FileSet) error {
	indexes := idx.Sort(file.Decls)
	mp := ast.NewCommentMap(set, file, file.Comments)
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading: %w", err)
	}
	var buffer bytes.Buffer
	i := 0

	for dstIndex, srcIndex := range indexes {
		dstDecl := file.Decls[dstIndex]
		dstStart, dstEnd := bounds(dstDecl, mp, set)
		prefix := src[i:dstStart]
		buffer.Write(prefix)

		srcDecl := file.Decls[srcIndex]
		srcStart, srcEnd := bounds(srcDecl, mp, set)
		infix := src[srcStart:srcEnd]
		buffer.Write(infix)

		i = dstEnd
	}

	suffix := src[i:]
	buffer.Write(suffix)

	dst := buffer.Bytes()
	if err := ioutil.WriteFile(path, dst, 0644); err != nil {
		return fmt.Errorf("failed writing: %w", err)
	}
	return nil
}

func arrangePackage(pkg *ast.Package, set *token.FileSet, filename string) error {
	docs := doc.New(pkg, "", doc.AllDecls|doc.PreserveAST)
	idx := index.New()

	appendValues(idx, docs.Consts)
	appendValues(idx, docs.Vars)
	appendFuncs(idx, docs.Funcs)

	for _, typ := range docs.Types {
		idx.Append(typ.Decl)
		appendValues(idx, typ.Consts)
		appendValues(idx, typ.Vars)
		appendFuncs(idx, typ.Funcs)
		appendFuncs(idx, typ.Methods)
	}

	for path, file := range pkg.Files {
		if filename == "" || filepath.Base(path) == filename {
			if err := arrangeFile(file, idx, path, set); err != nil {
				return fmt.Errorf("failed arranging file %q: %w", path, err)
			}
		}
	}
	return nil
}

func bounds(decl ast.Decl, mp ast.CommentMap, set *token.FileSet) (int, int) {
	start := decl.Pos()
	for _, group := range mp.Filter(decl).Comments() {
		if group.Pos() < start {
			start = group.Pos()
		}
	}

	end := decl.End()
	for _, group := range mp.Filter(decl).Comments() {
		if group.End() > end {
			end = group.End()
		}
	}

	return offset(start, set), offset(end, set)
}

func offset(pos token.Pos, set *token.FileSet) int {
	position := set.Position(pos)
	return position.Offset
}

func run(path string) error {
	set := token.NewFileSet()
	dir, filename, err := split(path)
	if err != nil {
		return fmt.Errorf("failed splitting path: %w", err)
	}

	packages, err := parser.ParseDir(set, dir, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed parsing: %w", err)
	}

	for _, pkg := range packages {
		if err := arrangePackage(pkg, set, filename); err != nil {
			return fmt.Errorf("failed arranging package %q: %w", pkg.Name, err)
		}
	}
	return nil
}

func split(path string) (string, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", fmt.Errorf("failed checking status of file: %w", err)
	}

	if info.IsDir() {
		return path, "", nil
	} else {
		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		return dir, filename, nil
	}
}
