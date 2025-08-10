package functions

import "github.com/sajuno/goon/language/lsp"

type GoToDefinitionInput struct {
	URI       string `json:"URI"`
	Line      int    `json:"Line"`
	Character int    `json:"Character"`
}

type GoToDefinitionOutput struct {
	Location *lsp.Location `json:"location"`
	Error    *lsp.Error    `json:"error"`
}
