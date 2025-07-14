package ingest

import (
	"fmt"
	"github.com/google/uuid"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Chunk holds (usually) blocks of code with semantic meaning in the context of an AI prompt
// They usually correspond directly to an AST node

type ChunkKind string

var (
	ChunkKindFunc ChunkKind = "func"
	ChunkKindType ChunkKind = "type"
	ChunkKindTest ChunkKind = "test"
)

type Chunk struct {
	// Stable UUID
	ID string

	// Content holds the entire block of code
	Content  string // The actual code block
	FilePath string
	Package  string
	Kind     ChunkKind

	// Name represents the name of the block of code. Usually var/type name
	Name string

	// Chunk position in file
	StartLine, EndLine int
}

// EmbeddedChunk also holds the chunk's vector
type EmbeddedChunk struct {
	Chunk
	Vector []float64
}

func findGoFiles(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// chunkFile reads a go file and deconstructs it into Chunks
func chunkFile(path string) ([]Chunk, error) {
	fset := token.NewFileSet()

	// get file level ast node
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", path, err)
	}

	// read in source file so we can the textual chunk content
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file on path %s: %w", path, err)
	}
	source := string(f)

	var chunks []Chunk
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			start := fset.Position(d.Pos())
			end := fset.Position(d.End())
			kind := ChunkKindFunc
			if strings.HasPrefix(d.Name.Name, "Test") {
				kind = ChunkKindTest
			}
			chunks = append(chunks, Chunk{
				ID:        uuid.NewString(),
				Content:   source[start.Offset:end.Offset],
				FilePath:  path,
				Package:   node.Name.Name,
				Kind:      kind,
				Name:      d.Name.Name,
				StartLine: start.Line,
				EndLine:   end.Line,
			})
		}
	}

	return chunks, nil
}
