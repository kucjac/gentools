package testdata

// GetPet returns one pet.
// @openapi method=GET route=/pets/{id} summary="Get a pet" response=200:Pet
func GetPet() {}

// ListPets returns all pets.
// @openapi method=GET route=/pets summary="List pets" response=200:Pet
func ListPets() {}
