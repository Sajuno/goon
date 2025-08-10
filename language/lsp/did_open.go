package lsp

import (
	"encoding/json"
)

func (c *Client) DidOpen(uri, langID, text string, version int) error {
	params := DidOpenTextDocumentParams{}
	params.TextDocument.URI = uri
	params.TextDocument.LanguageID = langID
	params.TextDocument.Version = version
	params.TextDocument.Text = text
	paramBytes, _ := json.Marshal(params)
	return c.send(&Message{
		Method: "textDocument/didOpen",
		Params: paramBytes,
	})
}
