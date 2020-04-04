package wrap

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestGetPets(t *testing.T) {
	store := &PetStoreService{
		pets: []*Pet{
			{
				Name:   "Lindy",
				Age:    7,
				Fluffy: true,
			},
		},
	}
	server := rpc.NewServer()
	defer server.Stop()

	err := server.RegisterReceiverWithName("store", store)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("can't listen:", err)
	}
	defer listener.Close()
	go server.ServeListener(listener)

	request := `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
	deadline := time.Now().Add(10 * time.Second)

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetDeadline(deadline)
	conn.Write([]byte(request))
	conn.(*net.TCPConn).CloseWrite()

	buf := []byte{}
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(buf[:n]))

}

