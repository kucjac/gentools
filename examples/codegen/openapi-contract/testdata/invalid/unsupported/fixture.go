package unsupported

type Pet struct {
	Values map[string]string
}

// Broken references an unsupported model shape.
// @openapi method=GET route=/x summary=bad response=200:Pet
func Broken() {}
