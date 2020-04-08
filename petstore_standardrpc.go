package petstore

/*
Types and methods in this file wrap the PetStore type in an interface
which can be handled by the standard lib rpc packages as expected.
*/
type PetStoreStandardRPCService struct {
	store *PetStore
}

type GetPetsArg string
type GetPetsRes []*Pet

// GetPets returns all the pets the petStoreService has.
func (s *PetStoreStandardRPCService) GetPets(arg GetPetsArg, response *GetPetsRes) error {
	res, err := s.store.GetPets()
	if err != nil {
		return err
	}

	response = (*GetPetsRes)(&res)
	return nil
}

type AddPetArg Pet
type AddPetRes Pet

// AddPet adds a pet to the petStoreService.
func (s *PetStoreStandardRPCService) AddPet(arg AddPetArg, response *AddPetRes) (err error) {
	err = s.store.AddPet(Pet(arg))
	if err != nil {
		return err
	}
	response = (*AddPetRes)(&arg)
	return nil
}

type RegisterPetOwnerArgs struct {
	OwnerName string
	Pet Pet
}

type RegisteredPetOwner struct {
	OwnerName string
	NewPet Pet
}
type RegisterPetOwnerRes struct {}

// RegisterPetOwner registers a pet to an owner.
func (s *PetStoreStandardRPCService) RegisterPetOwner(args RegisterPetOwnerArgs, res *RegisterPetOwnerRes) error {

	s.store.RegisterPetOwner(args.OwnerName, args.Pet)

	return nil
}