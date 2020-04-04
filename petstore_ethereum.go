package petstore

type Pet struct {
	Name         string
	Age          int
	Fluffy, Fast bool
}

type PetStoreEthereumService struct {
	pets []*Pet
}


// GetPets returns all the pets the store has.
func (s *PetStoreEthereumService) GetPets() ([]*Pet, error) {
	// Returns all pets.
	return s.pets, nil
}

// AddPet adds a pet to the store.
func (s *PetStoreEthereumService) AddPet(p Pet) error {
	if s.pets == nil {
		s.pets = []*Pet{}
	}
	s.pets = append(s.pets, &p)
	return nil
}

