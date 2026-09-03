module github.com/kucjac/gentools/internal/integration/scenarios/generated

go 1.27.0

require (
	github.com/kucjac/gentools v0.0.0-00010101000000-000000000000
	github.com/kucjac/gentools/internal/integration v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/kucjac/gentools => ../../../../

replace github.com/kucjac/gentools/internal/integration => ../..
