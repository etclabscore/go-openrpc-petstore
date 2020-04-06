package petstore

import (
	"encoding/json"
	"log"
	"reflect"

	openrpc_go_document "github.com/etclabscore/openrpc-go-document"
)

type PetStoreStdService struct {
	pets []*Pet
}

type None interface{}

type GetPetsArgs None
type GetPetsRes []*Pet

// GetPets returns all the pets the store has.
func (s *PetStoreStdService) GetPets(args GetPetsArgs, response *GetPetsRes) error {
	// Returns all pets.
	copy(*response, s.pets)
	return nil
}

type AddPetArg Pet
type AddPetRes Pet

// AddPet adds a pet to the store.
func (s *PetStoreStdService) AddPet(p AddPetArg, response *AddPetRes) (err error) {
	if s.pets == nil {
		s.pets = []*Pet{}
	}
	newPet := &Pet{}
	*newPet = Pet(p)
	s.pets = append(s.pets, newPet)
	petRes := AddPetRes(p)
	*response = petRes
	return nil
}

type StandardOpenRPCServiceProvider struct {
	service *PetStoreStdService
}

type StandardDiscoverArgs string
type StandardDiscoverRes map[string]interface{}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (s *PetStoreStdService) Discover(args StandardDiscoverArgs, response *StandardDiscoverRes) error {

	doc := openrpc_go_document.DocumentProvider(
		openrpc_go_document.DefaultGoRPCServiceProvider(s),
		openrpc_go_document.DefaultParseOptions(),
		)

	if doc == nil {
		log.Fatal("doc is nil")
	}

	err := doc.Discover()
	if err != nil {
		return err
	}

	b, err := json.Marshal(doc.Spec1())
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, response)
	if err != nil {
		return err
	}
	return nil
}
