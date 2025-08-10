package functions

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

type Definition struct {
	Name        string
	Description string
	Schema      map[string]any
}

func Definitions() []openai.FunctionDefinition {
	definitions := []openai.FunctionDefinition{
		{
			Name:        "did_open",
			Description: "Sends a notification that a text document has been opened (LSP's textDocument/didOpen)",
			Strict:      true,
			Parameters:  mustGenSchema(&DidOpenInput{}),
		},
		{
			Name:        "find_references",
			Description: "Find all references to a symbol at a given position in a document (LSP's textDocument/references)",
			Strict:      true,
			Parameters:  mustGenSchema(&FindReferencesInput{}),
		},
		{
			Name:        "go_to_definition",
			Description: "Find the definition of a symbol at a given position (LSP's goToDefinition)",
			Strict:      true,
			Parameters:  mustGenSchema(&GoToDefinitionInput{}),
		},
	}

	return definitions
}

func mustGenSchema(input any) map[string]any {
	reflector := jsonschema.Reflector{
		Anonymous:                 true,
		AllowAdditionalProperties: false,
		DoNotReference:            true,
		ExpandedStruct:            true,
	}

	r := reflector.Reflect(input)
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	var out map[string]any
	err = json.Unmarshal(b, &out)
	if err != nil {
		panic(err)
	}

	// jsonschema doesn't set this automatically, and it applies to all function definitions
	out["type"] = "object"

	return out
}
