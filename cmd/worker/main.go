package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/RGood/compute-go-vm/internal/generated/protos/echo"
	"google.golang.org/grpc"
)

var socketPath string = "/tmp/comms"

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (es *EchoServer) Ping(ctx context.Context, msg *echo.Message) (*echo.Message, error) {
	return msg, nil
}

func main() {
	socketAddr := fmt.Sprintf("%s/socket.sock", socketPath)

	println("Removing old socket")

	if err := os.Remove(socketAddr); err != nil {
		fmt.Printf("Error deleting old socket: %s\n", err.Error())
	}

	println("Listening to socket")

	lis, err := net.Listen("unix", socketAddr)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	println("Starting grpc server")

	server := grpc.NewServer()

	println("Registering echo")
	echo.RegisterEchoServer(server, &EchoServer{})

	println("Binding grpc server to socket")
	server.Serve(lis)
}
