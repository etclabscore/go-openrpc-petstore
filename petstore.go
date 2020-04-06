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

// GetPets returns all the pets the store has.
func (s *PetStore) GetPets() ([]*Pet, error) {
	// Returns all pets.
	return s.pets, nil
}

// AddPet adds a pet to the store.
func (s *PetStore) AddPet(p Pet) error {
	if s.pets == nil {
		s.pets = []*Pet{}
	}
	s.pets = append(s.pets, &p)
	return nil
}

func (s *PetStore) TryUpdatePet(ctx context.Context, p Pet) {
	for i, pet := range s.pets {
		if pet.Name == p.Name {
			*s.pets[i] = p
		}
	}
}

