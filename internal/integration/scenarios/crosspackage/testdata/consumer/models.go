package consumer

import "github.com/kucjac/gentools/internal/integration/scenarios/crosspackage/testdata/producer"

type Holder struct{ Box producer.Box[string] }
