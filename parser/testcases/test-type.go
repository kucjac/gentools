package testcases

import (
	"testing"
)

type testingType struct {
	embeddedType
	Integer int
}

func (t *testingType) PtrMethod(param string) error {
	return nil
}

func (t testingType) Method(paramNonPtr string) (err error) {
	return nil
}

type embeddedType struct {
	Imported *testing.T `gentools:"testme"`
}
