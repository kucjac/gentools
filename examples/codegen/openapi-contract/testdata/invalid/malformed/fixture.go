package malformed

// Broken has an unsupported method.
// @openapi method=TRACE route=/x summary=bad response=200:Pet
func Broken() {}
