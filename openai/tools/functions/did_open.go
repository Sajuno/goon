package functions

import "github.com/sajuno/goon/language/lsp"

type DidOpenInput struct {
	URI       string `json:"uri"`
	Text      string `json:"text"`
	LangID    string `json:"lang_id"`
	VersionID int    `json:"version_id"`
}

type DidOpenOutput struct {
	Error *lsp.Error `json:"error"`
}
