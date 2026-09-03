package duplicate

type Pet struct{}

// One documents the first operation.
// @openapi method=GET route=/x summary=one response=200:Pet
func One() {}

// Two conflicts with the first operation.
// @openapi method=GET route=/x summary=two response=200:Pet
func Two() {}
