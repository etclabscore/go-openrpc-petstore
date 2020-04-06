package petstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"testing"

	openrpc_go_document "github.com/etclabscore/openrpc-go-document"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

var standardStoreRPCService = &PetStoreStandardRPCService{
	store: &PetStore{pets:  []*Pet{
		{
			Name: "Bunny",
			Age: 64,
			Fluffy: true,
		},
	}},
}

type StandardOpenRPCServiceProvider struct {
	service *PetStoreStandardRPCService
}

type StandardDiscoverArgs string
type StandardDiscoverRes map[string]interface{}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (s *PetStoreStandardRPCService) Discover(args StandardDiscoverArgs, response *StandardDiscoverRes) error {

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

func TestRPCDocument_RPC(t *testing.T) {

	server := rpc.NewServer()
	err := server.RegisterName("petstore", standardStoreRPCService)
	if err != nil {
		t.Fatal(err)
	}

	// Set up a listener for our standard lib rpc server.
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}

	// Register a HTTP handler
	rpc.HandleHTTP()
	go func() {
		log.Printf("Serving RPC server on port %d", 1234)
		// Start accept incoming HTTP connections
		err = http.Serve(listener, server)
		if err != nil {
			log.Fatal("Error serving: ", err)
		}
	}()

	// Make connection to rpc server
	client, err := rpc.DialHTTP("tcp", listener.Addr().String())
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}
	defer client.Close()

	//make arguments object
	args := Pet{
		Name: "Chili",
		Age: 3,
		Fast: true,
	}
	// This will store returned result
	var result AddPetRes

	//call remote procedure with args
	err = client.Call("petstore.AddPet", args, &result)
	if err != nil {
		t.Fatal("call error", err)
	}

	//we got our result in result
	fmt.Printf("args=%v res=%v\n", args, result)

	allPetsJSON, _ := json.MarshalIndent(standardStoreRPCService.store.pets, "", "  ")
	fmt.Println("=> Added pet named" , result.Name, ", petstore pets state:", string(allPetsJSON))

	if len(standardStoreRPCService.store.pets) != 2 {
		t.Fatal("expect 2 pets")
	}

	r := make(map[string]interface{})
	discoverRes := StandardDiscoverRes(r)

	err = client.Call("petstore.Discover", "", &discoverRes)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(discoverRes, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))

	// Test that the reply marshals back into our datatype with field values' equivalence confirmed.
	check := &goopenrpcT.OpenRPCSpec1{}
	err = json.Unmarshal(b, check)
	if err != nil {
		t.Fatal(err)
	}

	// TOIMPROVE
	if len(check.Methods) != len([]string{"AddPet", "Discover", "GetPets"}) {
		t.Error(len(check.Methods))
	}
	if check.Methods[0].Name != "AddPet" {
		t.Error("want methods[0].Name = AddPet")
	}

}