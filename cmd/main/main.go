package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/RGood/compute-go-vm/internal/generated/protos/echo"
	"github.com/google/uuid"
	"github.com/weaveworks/footloose/pkg/cluster"
	"github.com/weaveworks/footloose/pkg/config"
	"google.golang.org/grpc"
)

func createMachine(c *cluster.Cluster, id string, backend string) *cluster.Machine {
	socketPath := fmt.Sprintf("/tmp/compute/%s", id)
	os.MkdirAll(socketPath, os.ModeDir)

	m := c.NewMachine(&config.Machine{
		Name:       fmt.Sprintf("worker-%s", id),
		Image:      "docker.io/library/compute:worker01",
		Privileged: true,
		PublicKey:  "machine-key",
		Backend:    backend,
		Volumes: []config.Volume{
			{
				Type:        "bind",
				Source:      socketPath,
				Destination: "/tmp/comms",
			},
		},
		Cmd: "/srv/server",
	})
	c.CreateMachine(m, 0)

	return m
}

var socketDir string = "/tmp/compute"

func dial(addr string, t time.Duration) (net.Conn, error) {
	return net.Dial("unix", addr)
}

func main() {
	// Create cluster under which to build containers
	c, err := cluster.New(config.Config{
		Cluster: config.Cluster{
			Name: "compute",
		},
	})
	if err != nil {
		panic(err)
	}

	// Create a key store (required for creating containers)
	s := cluster.NewKeyStore("/tmp/keystore")
	s.Init()

	// Store a dummy public key
	err = s.Store("machine-key", "machine-key")

	// Set the container's key store to the one we made
	c.SetKeyStore(s)

	// Create a custom machine ID
	id := uuid.NewString()

	// Create machine
	m := createMachine(c, id, "docker")
	defer c.DeleteMachine(m, 0)

	// Connect to the gRPC socket created by the machine
	conn, err := grpc.Dial(fmt.Sprintf("/tmp/compute/%s/socket.sock", id), grpc.WithInsecure(), grpc.WithDialer(dial))
	if err != nil {
		println(err.Error())
		return
	}

	// Instantiate an echo client using that connection
	echoClient := echo.NewEchoClient(conn)

	// Call the echo client N times and verify it succeeded
	for i := 0; i < 10; i++ {
		res, err := echoClient.Ping(context.Background(), &echo.Message{
			Message: fmt.Sprintf("foo: %d", i),
		})
		if err != nil {
			println(err.Error())
			return
		}
		println(res.Message)
	}

	// Create vm config from template
	// Instantiate vm

	// // Message it N times via the socket

	// Teardown
}
