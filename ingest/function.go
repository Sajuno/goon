package ingest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Function struct {
	Source  string
	Path    string
	Package string
}

type FindFunctionQuery struct {
	Name    string
	Package string
}

func FindFunction(rootPath string, q FindFunctionQuery) (Function, error) {
	var source string
	var filePath string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return nil
		}

		if node.Name.Name != "" && node.Name.Name != q.Package {
			return nil
		}

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == q.Name {
				start := fset.Position(fn.Pos()).Offset
				end := fset.Position(fn.End()).Offset

				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				source = string(data[start:end])
				filePath = path
				return io.EOF // break walk early
			}
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return Function{}, err
	}

	if source == "" {
		return Function{}, fmt.Errorf("function %s not found", q.Name)
	}

	return Function{Source: source, Path: filePath, Package: q.Package}, nil
}
