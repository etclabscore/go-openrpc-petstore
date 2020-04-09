package petstore

import (
	"net"
	"testing"

	openRPCDoc "github.com/etclabscore/go-openrpc-reflect"
	"github.com/ethereum/go-ethereum/rpc"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

const maxReadSize = 1024 * 1024

var petStoreService = &PetStore{pets: []*Pet{
	{
		Name:   "Fluffy",
		Age:    14,
		Fluffy: true,
		Fast:   false,
	},
}}

func TestRPCDocument_EthereumRPC(t *testing.T) {
	server := rpc.NewServer()
	defer server.Stop()

	// Register the petStoreService service with the server, just like usual
	// for go-ethereum/rpc
	err := server.RegisterName("petStoreService", petStoreService)
	if err != nil {
		t.Fatal(err)
	}

	// Server settings are top-level; one ~~server~~ API, one document.
	serverConfigurationP := openRPCDoc.DefaultServerServiceProvider
	serverConfigurationP.ServiceOpenRPCInfoFn = func() goopenrpcT.Info {
		return goopenrpcT.Info{
			Title:          "My Ethereum-Style Petstore Service",
			Description:    "Oooohhh!",
			TermsOfService: "https://google.com/rtfm",
			Contact:        goopenrpcT.Contact{},
			License:        goopenrpcT.License{},
			Version:        "v0.0.0-beta",
		}
	}

	// Get our document provider from the serviceProvider.
	serverDoc := openRPCDoc.NewReflectDocument(serverConfigurationP)

	rp := openRPCDoc.EthereumRPCDescriptor

	serverDoc.Reflector.RegisterReceiver(petStoreService, rp)

	server.RegisterName("rpc", serverDoc)


	// Set up a listener for our go-ethereum/rpc server.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener) // Run it.

	// Send some test requests.
	requests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n",
		`{"jsonrpc":"2.0","id":1,"method":"rpc_discover"}` + "\n",
	}

	for _, request := range requests {
		makeRequest(t, request, listener)
	}


}
