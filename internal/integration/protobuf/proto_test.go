package protobuf

import (
	"testing"

	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
	"github.com/stretchr/testify/require"
)

func TestProtobufs(t *testing.T) {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{"."}})
	require.NoError(t, err)

	this := pkgs.MustGetByPath("github.com/kucjac/gentools/internal/integration/protobuf")
	tm, ok := this.GetType("TestingMessage")
	require.True(t, ok)

	st, ok := tm.(*types.Struct)
	require.True(t, ok)

	for _, field := range st.Fields {
		if field.Name == "Any" {
		}
	}
}
