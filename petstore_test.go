package petstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
	bufPretty, err := json.MarshalIndent(pretty["result"], "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	// Write it cough hack cough
	if strings.Contains(request, "rpc_discover") {
		err = ioutil.WriteFile(filepath.Join("generated", "ethereum.json"), bufPretty, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("<--", string(bufPretty))
	fmt.Println()
}

