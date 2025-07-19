package golang

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ChunkKind string

func (k ChunkKind) String() string {
	return string(k)
}

var (
	ChunkKindFunc       ChunkKind = "func"
	ChunkKindTest       ChunkKind = "test"
	ChunkKindMethod     ChunkKind = "method"
	ChunkKindStruct     ChunkKind = "struct"
	ChunkKindInterface  ChunkKind = "interface"
	ChunkKindTypeAlias  ChunkKind = "type_alias"
	ChunkKindConstBlock ChunkKind = "const"
	ChunkKindVarBlock   ChunkKind = "var"
	ChunkKindUnknown    ChunkKind = "unknown"
)

// Chunk holds (usually) blocks of code with semantic meaning in the context of an AI prompt
// They correspond to an AST node or otherwise have semantic meaning
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

	// Optional comments
	Doc string

	// Receiver is only populated for methods
	ReceiverName string
}

// Sha256 returns the chunks checksum based on a couple of fields
func (c Chunk) Sha256() string {
	h := sha256.New()
	_, _ = io.WriteString(h, c.FilePath)
	_, _ = io.WriteString(h, c.Package)
	_, _ = io.WriteString(h, strconv.Itoa(c.StartLine))
	_, _ = io.WriteString(h, strconv.Itoa(c.EndLine))
	_, _ = io.WriteString(h, c.Content)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func ChunkRepository(path string) ([]Chunk, error) {
	files, err := findFiles(path)
	if err != nil {
		return nil, err
	}

	var all []Chunk
	for _, f := range files {
		chunks, err := chunkFile(f)
		if err != nil {
			log.Printf("❌ Error parsing %s: %v", f, err)
			continue
		}
		all = append(all, chunks...)
	}
	return all, nil
}

func findFiles(path string) ([]string, error) {
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
			receiver, ok := getReceiver(d)
			if ok {
				kind = ChunkKindMethod
			}
			chunks = append(chunks, Chunk{
				ID:           uuid.NewString(),
				Content:      source[start.Offset:end.Offset],
				FilePath:     path,
				Package:      node.Name.Name,
				Kind:         kind,
				Name:         d.Name.Name,
				StartLine:    start.Line,
				EndLine:      end.Line,
				Doc:          d.Doc.Text(),
				ReceiverName: receiver,
			})
		case *ast.GenDecl:
			// Skip import tokens
			if d.Tok == token.IMPORT {
				continue
			}

			for _, spec := range d.Specs {
				var name string

				start := fset.Position(spec.Pos())
				end := fset.Position(spec.End())
				kind := classifyGenDecl(d)

				switch s := spec.(type) {
				case *ast.TypeSpec:
					name = s.Name.Name
				case *ast.ValueSpec:
					if len(s.Names) > 0 {
						name = s.Names[0].Name
					}
				default:
					name = "" // explicitly unnamed
				}

				chunks = append(chunks, Chunk{
					ID:        uuid.NewString(),
					Content:   source[start.Offset:end.Offset],
					FilePath:  path,
					Package:   node.Name.Name,
					Kind:      kind,
					Name:      name,
					StartLine: start.Line,
					EndLine:   end.Line,
					Doc:       d.Doc.Text(),
				})
			}
		}
	}

	return chunks, nil
}

func classifyGenDecl(decl *ast.GenDecl) ChunkKind {
	switch decl.Tok {
	case token.TYPE:
		for _, spec := range decl.Specs {
			s, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			switch s.Type.(type) {
			case *ast.StructType:
				return ChunkKindStruct
			case *ast.InterfaceType:
				return ChunkKindInterface
			default:
				return ChunkKindTypeAlias // no other options
			}
		}
	case token.VAR:
		return ChunkKindVarBlock
	case token.CONST:
		return ChunkKindConstBlock
	default:
	}
	return ChunkKindUnknown
}

// getReceiver determines the receiver name and if it exists
func getReceiver(fn *ast.FuncDecl) (string, bool) {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", false
	}

	recvType := fn.Recv.List[0].Type

	switch expr := recvType.(type) {
	case *ast.StarExpr:
		if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name, true
		}
	case *ast.Ident:
		return expr.Name, true
	}

	return "", false
}
