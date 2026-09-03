package unknowndirective

type Pet struct{}

// Broken contains an unknown directive key.
// @openapi method=GET route=/x summary=bad response=200:Pet extra=value
func Broken() {}
