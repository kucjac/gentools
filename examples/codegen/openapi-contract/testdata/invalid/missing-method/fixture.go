package missingmethod

type Pet struct{}

// Broken omits the method.
// @openapi route=/x summary=bad response=200:Pet
func Broken() {}
