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
	Source string
	Path   string
}

func FindFunction(rootPath, name string) (Function, error) {
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

		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == name {
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
		return Function{}, fmt.Errorf("function %s not found", name)
	}

	return Function{Source: source, Path: filePath}, nil
}
