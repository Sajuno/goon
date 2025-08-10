package lsp

import (
	"encoding/json"
	"fmt"
)

func (c *Client) FindReferences(uri string, line, char int) ([]Location, error) {
	params := TextDocumentPositionParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Position:     Position{Line: line, Character: char},
	}
	paramBytes, _ := json.Marshal(params)
	msg := newMessage(
		"textDocument/references",
		paramBytes,
	)
	if err := c.send(msg); err != nil {
		return nil, err
	}
	resp, err := c.read()
	if err != nil {
		return nil, err
	}
	var locations []Location
	if err := json.Unmarshal(resp.Result, &locations); err != nil {
		var single Location
		if err2 := json.Unmarshal(resp.Result, &single); err2 == nil {
			locations = append(locations, single)
		} else {
			return nil, fmt.Errorf("failed to decode references: %w", err)
		}
	}
	return locations, nil
}
