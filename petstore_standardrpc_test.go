package petstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc" // <--  Using the standard library RPC package.
	"os"
	"path/filepath"
	"testing"

	openRPCDoc "github.com/etclabscore/openrpc-go-document"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

// DocService is a very thin wrapper around the Document object,
// which we'll use to build a standard library compatible interface
// returning the document with an arbitrary type (here, json.RawMessage).
type DocService struct {
	*openRPCDoc.Document
}

type DiscoverArg string
type DiscoverReply json.RawMessage

// Discover wraps Document.Discover in a method whose signature fulfills
// the conventions of go standard library rpc server registration.
func (d *DocService) Discover(arg DiscoverArg, reply *DiscoverReply) error {
	got, err := d.Document.Discover()
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(got, "", "  ")

	*reply = out
	return err
}

func TestRPCDocument_StandardRPC(t *testing.T) {

	server := rpc.NewServer()

	standardStoreRPCService := &PetStoreStandardRPCService{

		// Populated with some fake data.
		store: &PetStore{pets: []*Pet{
			{
				Name:   "Bunny",
				Age:    64,
				Fluffy: true,
			},
		}},
	}

	// Register our service.
	err := server.RegisterName("PetStore", standardStoreRPCService)
	if err != nil {
		t.Fatal(err)
	}

	// Server settings are top-level; one ~~server~~ API, one document.
	//
	// Get the server (not service!) provider default.
	serverConfigurationP := openRPCDoc.DefaultServerServiceProvider

	// Modify the server config.
	serverConfigurationP.ServiceOpenRPCInfoFn = func() goopenrpcT.Info {
		return goopenrpcT.Info{
			Title:          "My Standard Service",
			Description:    "Aaaaahh!",
			TermsOfService: "https://google.com/rtfm",
			Contact:        goopenrpcT.Contact{},
			License:        goopenrpcT.License{},
			Version:        "v0.0.0-beta",
		}
	}

	// Create a new "reflectable" document for this server.
	serverDoc := openRPCDoc.NewReflectDocument(serverConfigurationP)

	// Get the service provider default for our RPC API style,
	// in this case, the Go standard lib.
	sp := openRPCDoc.StandardRPCDescriptor

	// Register our receiver-based service standardService.
	serverDoc.Reflector.RegisterReceiverWithName("PetStore", standardStoreRPCService, sp)

	wrapped := &DocService{serverDoc}

	err = server.RegisterName("rpc", wrapped)
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

	// Do some test RPC calls against the API.
	args := Pet{
		Name: "Chili",
		Age:  3,
		Fast: true,
	}
	var result AddPetRes
	err = client.Call("PetStore.AddPet", args, &result)
	if err != nil {
		t.Fatal("call error", err)
	}

	//we got our result in result
	t.Logf("args=%v res=%v\n", args, result)

	allPetsJSON, _ := json.MarshalIndent(standardStoreRPCService.store.pets, "", "  ")
	t.Log("=> Added pet named", result.Name, ", petstore pets state:", string(allPetsJSON))

	// Check that a pet was actually added in mem.
	if len(standardStoreRPCService.store.pets) != 2 {
		t.Fatal("expect 2 pets")
	}

	// Call the RPC Discover method.
	dreply := json.RawMessage{}
	err = client.Call("rpc.Discover", "none", &dreply)
	if err != nil {
		t.Fatal(err)
	}

	// Log for visibility
	t.Log(string(dreply))
	err = ioutil.WriteFile(filepath.Join("generated", "standard.json"), dreply, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// Test that the reply marshals back into our datatype with field values' equivalence confirmed.
	check := &goopenrpcT.OpenRPCSpec1{}
	err = json.Unmarshal(dreply, check)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Deeper and more complete tests that the OpenRPC API description document
	// contains the elements we expect.
	if len(check.Methods) != len([]string{"AddPet", "Discover", "GetPets"}) {
		t.Error(len(check.Methods))
	}
	if n := check.Methods[0].Name; n != "PetStore.AddPet" {
		t.Error("want methods[0].Name = AddPet, got", n)
	}

}
