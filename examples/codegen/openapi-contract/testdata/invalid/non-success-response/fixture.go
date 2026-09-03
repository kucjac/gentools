package nonsuccessresponse

type Pet struct{}

// Broken declares an error response in the success-only example scope.
// @openapi method=GET route=/x summary=bad response=404:Pet
func Broken() {}
