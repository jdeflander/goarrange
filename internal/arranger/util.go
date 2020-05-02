package arranger

import (
	"go/ast"
	"go/token"
)

func bounds(decls []ast.Decl, index int, mp ast.CommentMap, set *token.FileSet) (int, int) {
	decl := decls[index]
	minStart := minStart(decls, index, mp)
	start := decl.Pos()
	for _, group := range mp.Filter(decl).Comments() {
		if group.Pos() > minStart && group.Pos() < start {
			start = group.Pos()
		}
	}
	end := end(decl, mp)
	return offset(start, set), offset(end, set)
}

func end(decl ast.Decl, mp ast.CommentMap) token.Pos {
	end := decl.End()
	for _, group := range mp.Filter(decl).Comments() {
		if group.End() > end {
			end = group.End()
		}
	}
	return end
}

func minStart(decls []ast.Decl, index int, mp ast.CommentMap) token.Pos {
	if index == 0 {
		return token.NoPos
	}
	decl := decls[index-1]
	return end(decl, mp)
}

func offset(pos token.Pos, set *token.FileSet) int {
	position := set.Position(pos)
	return position.Offset
}
