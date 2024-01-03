package audiogenerator

import "context"

type AudioGeneratorDoer interface {
	Generate(ctx context.Context, language, text string) (*Result, error)
}

const (
	Language_pt_BR = "pt-BR"
	Language_en_US = "en-US"
)

type Result struct {
	URL string `json:"url"`
}
