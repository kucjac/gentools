package missingroute

type Pet struct{}

// Broken omits the route.
// @openapi method=GET summary=bad response=200:Pet
func Broken() {}
