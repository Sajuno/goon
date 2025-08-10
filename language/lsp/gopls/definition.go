package gopls

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sajuno/goon/language/lsp"
)

func (c *Client) GoToDefinition(uri string, line, char int) (*lsp.Location, error) {
	params := lsp.TextDocumentPositionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Position:     lsp.Position{Line: line, Character: char},
	}
	paramBytes, _ := json.Marshal(params)
	msg := &lsp.Message{
		ID:     uuid.NewString(),
		Method: "textDocument/definition",
		Params: paramBytes,
	}

	if err := c.send(msg); err != nil {
		return nil, err
	}

	resp, err := c.read()
	if err != nil {
		return nil, err
	}

	var locations []lsp.Location
	if err := json.Unmarshal(resp.Result, &locations); err != nil {
		var single lsp.Location
		if err2 := json.Unmarshal(resp.Result, &single); err2 == nil {
			locations = append(locations, single)
		} else {
			return nil, fmt.Errorf("failed to decode location: %w", err)
		}
	}
	if len(locations) == 0 {
		return nil, nil
	}

	return &locations[0], nil
}
