package missingresponse

type Pet struct{}

// Broken omits the response.
// @openapi method=GET route=/x summary=bad
func Broken() {}
