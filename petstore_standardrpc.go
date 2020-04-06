package petstore

/*
Types and methods in this file wrap the PetStore type in an interface
which can be handled by the standard lib rpc packages as expected.
*/
type PetStoreStandardRPCService struct {
	store *PetStore
}

type GetPetsArgs string
type GetPetsRes []*Pet

// GetPets returns all the pets the store has.
func (s *PetStoreStandardRPCService) GetPets(args GetPetsArgs, response *GetPetsRes) error {
	res, err := s.store.GetPets()
	if err != nil {
		return err
	}

	response = (*GetPetsRes)(&res)
	return nil
}

type AddPetArg Pet
type AddPetRes Pet

// AddPet adds a pet to the store.
func (s *PetStoreStandardRPCService) AddPet(p AddPetArg, response *AddPetRes) (err error) {
	err = s.store.AddPet(Pet(p))
	if err != nil {
		return err
	}
	response = (*AddPetRes)(&p)
	return nil
}
