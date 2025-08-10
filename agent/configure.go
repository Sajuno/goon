package agent

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/openai/tools/functions"
	"github.com/sashabaranov/go-openai"
)

const (
	codeAnalysisInstructions = `
You are a code analysis assistant with access to Language Server Protocol (LSP) tools. You can inspect a Go codebase by issuing structured tool calls like:

- did_open_file(uri, lang_id, text, version)
- go_to_definition(uri, line, character)
- find_references(uri, line, character)

You are not given some context up front, but you must explore the code by reading files and issuing tool requests. Assume the user’s question refers to real symbols in the codebase. Ask for files to be opened or symbols to be resolved before answering.

Always reason step-by-step:
1. Figure out what symbol or concept the user is asking about.
2. Determine what file or location you want to inspect.
3. Use tools to look up definitions, references, types, or related code.
4. Accumulate what you’ve learned before giving an answer.

Prefer concise and structured explanations.
Be blunt if something is poorly written or ambiguous.
If you can't find enough context to answer, say so and explain what info is missing.

You are not a search engine. You are a senior dev with tools — use them to figure things out.
`
)

func (a *Agent) Configure(ctx context.Context) error {
	var tools []openai.AssistantTool
	for _, def := range functions.Definitions() {
		tools = append(tools, openai.AssistantTool{
			Type:     openai.AssistantToolTypeFunction,
			Function: &def,
		})
	}

	_, err := a.openai.ModifyAssistant(ctx, a.cfg.ID, openai.AssistantRequest{
		Name:         ptr("Goon"),
		Description:  ptr("Goon's code analysis assistent"),
		Instructions: ptr(codeAnalysisInstructions),
		Tools:        tools,
	})
	if err != nil {
		return fmt.Errorf("failed to update openai assistant: %w", err)
	}

	return nil
}

func ptr[T any](v T) *T {
	return &v
}
