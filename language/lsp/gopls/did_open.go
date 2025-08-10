package gopls

import (
	"encoding/json"
	"github.com/sajuno/goon/language/lsp"
)

func (c *Client) DidOpen(uri, langID, text string, version int) error {
	params := lsp.DidOpenTextDocumentParams{}
	params.TextDocument.URI = uri
	params.TextDocument.LanguageID = langID
	params.TextDocument.Version = version
	params.TextDocument.Text = text
	paramBytes, _ := json.Marshal(params)
	return c.send(&lsp.Message{
		Method: "textDocument/didOpen",
		Params: paramBytes,
	})
}
