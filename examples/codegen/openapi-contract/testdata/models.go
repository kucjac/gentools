package testdata

import (
	_ "github.com/kucjac/gentools/examples/codegen/openapi-contract/testdata/other"
	"github.com/kucjac/gentools/examples/codegen/openapi-contract/testdata/shared"
)

// Pet is the shared public response model.
type Pet struct {
	// ID identifies the pet.
	ID int `json:"id"`
	// Name is the display name.
	Name    string          `json:"name"`
	Tags    []string        `json:"tags,omitempty"`
	Rating  float32         `json:"rating"`
	Score   float64         `json:"score"`
	Owner   *shared.Profile `json:"owner,omitempty"`
	Ignored string          `json:"-"`
}
