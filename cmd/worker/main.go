package main

import (
	"context"
	"fmt"
	"net"

	"github.com/RGood/compute-go-vm/internal/generated/protos/echo"
	"google.golang.org/grpc"
)

var socketPath string = "/tmp/comms"

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (es *EchoServer) Ping(ctx context.Context, msg *echo.Message) (*echo.Message, error) {
	msg.Message = "pong: " + msg.Message
	return msg, nil
}

func main() {
	lis, err := net.Listen("unix", fmt.Sprintf("%s/socket.sock", socketPath))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	echo.RegisterEchoServer(server, &EchoServer{})

	server.Serve(lis)
}
