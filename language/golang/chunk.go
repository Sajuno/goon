package golang

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
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

func (c Chunk) IsInvokable() bool {
	return c.Kind == ChunkKindMethod || c.Kind == ChunkKindFunc
}

func ChunkRepository(path string) ([]Chunk, error) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Dir:  path,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	var allChunks []Chunk

	for _, pkg := range pkgs {
		fset := pkg.Fset
		info := pkg.TypesInfo

		// Collect all chunks and build object -> FQN map
		var pkgChunks []Chunk
		for _, file := range pkg.Syntax {
			chunks, err := chunkASTFile(file, fset, pkg.PkgPath, info)
			if err != nil {
				return nil, fmt.Errorf("failed to chunk file %s: %w", file.Name.Name, err)
			}

			for _, chunk := range chunks {
				pkgChunks = append(pkgChunks, chunk)
			}
		}

		allChunks = append(allChunks, pkgChunks...)
	}

	return allChunks, nil
}

// chunkFile reads a go file and deconstructs it into Chunks
func chunkASTFile(file *ast.File, fset *token.FileSet, pkgPath string, info *types.Info) ([]Chunk, error) {
	var chunks []Chunk

	filename := fset.Position(file.Pos()).Filename
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}
	source := string(b)

	for _, decl := range file.Decls {
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
				FilePath:  filename,
				Package:   pkgPath,
				Kind:      kind,
				Name:      d.Name.Name,
				StartLine: start.Line,
				EndLine:   end.Line,
				Doc:       d.Doc.Text(),
			})

		case *ast.GenDecl:
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
					FilePath:  filename,
					Package:   pkgPath,
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
