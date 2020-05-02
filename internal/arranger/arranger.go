package arranger

import (
	"go/ast"
	"go/token"

	"github.com/jdeflander/goarrange/internal/arranger/internal/index"
)

// Arranger represents a source code arranger for go packages.
type Arranger struct {
	index index.Index
	set   *token.FileSet
}

// New creates a new arranger for the given package and file set.
//
// The given package must have been parsed with the given file set.
func New(pkg *ast.Package, set *token.FileSet) Arranger {
	idx := index.New(pkg)
	return Arranger{
		index: idx,
		set:   set,
	}
}

// Arrange arranges the given file with the given arranger.
//
// The given file must be part of the given arranger's package, and its contents should be represented by the given
// bytes. This method returns an arranged copy of the given contents.
func (a Arranger) Arrange(file *ast.File, src []byte) []byte {
	indexes := a.index.Sort(file.Decls)
	size := len(src)
	dst := make([]byte, size)
	dstOffset := 0
	mp := ast.NewCommentMap(a.set, file, file.Comments)
	srcOffset := 0

	for dstIndex, srcIndex := range indexes {
		dstPrefix := dst[dstOffset:]
		dstStart, dstEnd := bounds(file.Decls, dstIndex, mp, a.set)
		srcPrefix := src[srcOffset:dstStart]
		dstOffset += copy(dstPrefix, srcPrefix)

		dstInfix := dst[dstOffset:]
		srcStart, srcEnd := bounds(file.Decls, srcIndex, mp, a.set)
		srcInfix := src[srcStart:srcEnd]
		dstOffset += copy(dstInfix, srcInfix)

		srcOffset = dstEnd
	}

	dstSuffix := dst[dstOffset:]
	srcSuffix := src[srcOffset:]
	copy(dstSuffix, srcSuffix)
	return dst
}

// Arranged checks whether the given file is arranged according to the given arranger.
//
// The given file must be part of the given arranger's package.
func (a Arranger) Arranged(file *ast.File) bool {
	return a.index.Sorted(file.Decls)
}
