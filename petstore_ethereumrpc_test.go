package petstore

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	openRPCDoc "github.com/etclabscore/openrpc-go-document"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-openapi/spec"
	goopenrpcT "github.com/gregdhill/go-openrpc/types"
)

const maxReadSize = 1024 * 1024

var store = &PetStore{pets: []*Pet{
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

	// Register the store service with the server, just like usual
	// for go-ethereum/rpc
	err := server.RegisterName("store", store)
	if err != nil {
		t.Fatal(err)
	}

	// Get a Spec1 service type wrapped around the server.
	doc := openRPCDoc.DocumentProvider(
		openRPCDoc.DefaultEthereumServiceProvider(store),
		openRPCDoc.DefaultEthereumParseOptions(),
	)

	doc.Discover()

	// Register the DocumentService as a service receiver.
	err = server.RegisterName("rpc", openRPCDoc.NewDiscoverService(doc))
	if err != nil {
		t.Fatal(err)
	}

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

/*

 */

// MySpecialOpenRPCService is a wrapper service that can be used to extend existing metadata services.
type MySpecialOpenRPCService struct {
	openRPCDocument *openRPCDoc.Document
}

func (m *MySpecialOpenRPCService) MySpecialDiscover() (res map[string]interface{}, err error) {
	err = m.openRPCDocument.Discover()
	if err != nil {
		return nil, err
	}
	marshaled, err := json.MarshalIndent(m.openRPCDocument.Spec1(), "", "  ")
	if err != nil {
		return nil, err
	}

	res = make(map[string]interface{})
	err = json.Unmarshal(marshaled, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *MySpecialOpenRPCService) MySpecialGetMethodNames() ([]string, error) {
	out := []string{}
	for _, m := range m.openRPCDocument.Spec1().Methods {
		out = append(out, m.Name)
	}
	return out, nil
}

func TestRPCDocument_EthereumRPCExtendedService(t *testing.T) {
	server := rpc.NewServer()
	defer server.Stop()

	// Register the store service with the server, just like usual
	// for go-ethereum/rpc
	err := server.RegisterName("store", store)
	if err != nil {
		t.Fatal(err)
	}

	// Get a Spec1 service type wrapped around the server.
	doc := openRPCDoc.DocumentProvider(
		&openRPCDoc.ServerProviderService{
			ServiceCallbacks:           openRPCDoc.DefaultServiceCallbacksEthereum(store),
			ServiceOpenRPCInfo:         func() goopenrpcT.Info { return goopenrpcT.Info{} },
			ServiceOpenRPCExternalDocs: func() *goopenrpcT.ExternalDocs { return nil },
		},
		&openRPCDoc.DocumentProviderParseOpts{

			SchemaMutationFns: []func(*spec.Schema) error{
				openRPCDoc.SchemaMutationExpand,
				openRPCDoc.SchemaMutationRemoveDefinitionsField,
			},
			MethodBlackList: []string{},
		})

	// Or, use a different service pattern to configure the endpoint using the OpenRPC Spec1 instance.
	myService := &MySpecialOpenRPCService{openRPCDocument: doc}
	err = server.RegisterName("special", myService)
	if err != nil {
		t.Fatal(err)
	}

	// Set up a listener for our go-ethereum/rpc server.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener) // Run it.

	// Send some test requests.
	requests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"special_mySpecialDiscover"}` + "\n",
		`{"jsonrpc":"2.0","id":1,"method":"special_mySpecialGetMethodNames"}` + "\n",
	}

	for _, request := range requests {
		makeRequest(t, request, listener)
	}
}

func makeRequest(t *testing.T, request string, listener net.Listener) {
	fmt.Println("-->", request)
	deadline := time.Now().Add(10 * time.Second)
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetDeadline(deadline)
	conn.Write([]byte(request))
	conn.(*net.TCPConn).CloseWrite()

	buf := make([]byte, maxReadSize)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	pretty := make(map[string]interface{})
	err = json.Unmarshal(buf[:n], &pretty)
	if err != nil {
		t.Fatal(err)
	}
	if e, ok := pretty["error"]; ok {
		t.Fatal(e)
	}
	bufPretty, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("<--", string(bufPretty))
	fmt.Println()
}
