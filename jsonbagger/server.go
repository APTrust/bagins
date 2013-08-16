package main

import (
	"fmt"
	"github.com/APTrust/bagins/jsonbagger"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"path/filepath"
)

const (
	BAGDIR = filepath.Join(os.TempDir(), "rpcbags/")
	PORT   = ":8222"
)

func startServer() {
	bagger := jsonbagger.NewJSONBagger()

	server := rpc.NewServer()
	server.Register(bagger)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error Listening:", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}
	}

	go server.ServerCodec(jsonrpc.NewServerCodec(conn))
}
