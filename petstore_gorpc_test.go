package petstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"
)

var standardStore = &PetStoreStdService{
	pets: []*Pet{
		{
			Name: "Bunny",
			Age: 64,
			Fluffy: true,
		},
	},
}

func TestRPCDocument_RPC(t *testing.T) {

	server := rpc.NewServer()
	err := server.RegisterName("petstore", standardStore)
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

}