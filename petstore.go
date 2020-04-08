package petstore

import "context"

type Pet struct {
	Name         string
	Age          int
	Fluffy, Fast bool
}

type PetStore struct {
	pets []*Pet
}

// petstore_getPets

// GetPets returns all the pets the petStoreService has.
func (s *PetStore) GetPets() ([]*Pet, error) {
	// Returns all pets.
	return s.pets, nil
}

// AddPet adds a pet to the petStoreService.
func (s *PetStore) AddPet(newPet Pet) error {
	if s.pets == nil {
		s.pets = []*Pet{}
	}
	s.pets = append(s.pets, &newPet)
	return nil
}

// TryUpdatePet has context.Context as the first param field, which should be ignored.
// It does not have a return value.
func (s *PetStore) TryUpdatePet(ctx context.Context, pet Pet) {
	for i, pet := range s.pets {
		if pet.Name == pet.Name {
			*s.pets[i] = *pet
		}
	}
}

// RegisterPetOwner registers a pet owner.
func (s *PetStore) RegisterPetOwner(ownerName string, pet Pet) (response Pet, err error) {
	return Pet{}, nil
}

func (s *PetStore) private() {

}