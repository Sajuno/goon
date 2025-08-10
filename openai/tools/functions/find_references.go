package functions

import "github.com/sajuno/goon/language/lsp"

type FindReferencesInput struct {
	URI       string `json:"uri"`
	Line      int    `json:"line"`
	Character int    `json:"character"`
}

type FindReferencesOutput struct {
	Locations []lsp.Location `json:"locations"`
	Error     *lsp.Error     `json:"error,omitempty"`
}
