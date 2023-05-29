package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RGood/compute-go-vm/internal/generated/protos/echo"
	"github.com/google/uuid"
	"github.com/weaveworks/footloose/pkg/cluster"
	"github.com/weaveworks/footloose/pkg/config"
	"google.golang.org/grpc"
)

func createMachine(c *cluster.Cluster, id string, backend string) (*cluster.Machine, error) {
	socketPath := fmt.Sprintf("/tmp/compute/%s", id)
	err := os.MkdirAll(socketPath, 0777)
	if err != nil {
		return nil, err
	}

	m := c.NewMachine(&config.Machine{
		Name:       fmt.Sprintf("worker-%s", id),
		Image:      "docker.io/library/compute:worker-js",
		Privileged: false,
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
	err = c.CreateMachine(m, 0)
	if err != nil {
		return nil, err
	}

	return m, nil
}

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
	s.Store("machine-key", "machine-key")

	// Set the container's key store to the one we made
	c.SetKeyStore(s)

	// Create a custom machine ID
	id := uuid.NewString()

	// Create machine
	mStart := time.Now()
	m, err := createMachine(c, id, "docker")
	if err != nil {
		log.Fatalf("Error creating machine: %s\n", err.Error())
	}
	mDuration := time.Since(mStart)
	fmt.Printf("Machine started in %s\n", mDuration)
	defer c.DeleteMachine(m, 0)

	// Connect to the gRPC socket created by the machine
	conn, err := grpc.Dial(fmt.Sprintf("/tmp/compute/%s/socket.sock", id), grpc.WithInsecure(), grpc.WithAuthority("localhost"), grpc.WithDialer(dial))
	if err != nil {
		println(err.Error())
		return
	}

	// Instantiate an echo client using that connection
	echoClient := echo.NewEchoClient(conn)

	requests := 1000

	// Call the echo client N times and verify it succeeded
	start := time.Now()
	for i := 0; i < requests; i++ {
		_, err := echoClient.Ping(context.Background(), &echo.Message{
			Message: fmt.Sprintf("foo: %d", i),
		})
		if err != nil {
			println(err.Error())
			return
		}
	}
	d := time.Since(start)

	fmt.Printf("Sync: %d requests made in: %s\n", requests, d)
	fmt.Printf("Avg. Req. Duration: %s\n", d/time.Duration(requests))

	wg := sync.WaitGroup{}
	start = time.Now()
	total := atomic.Int32{}
	// Message it N times via the socket
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			s := time.Now()
			_, err := echoClient.Ping(context.Background(), &echo.Message{
				Message: fmt.Sprintf("foo: %d", i),
			})
			if err != nil {
				println(err.Error())
				return
			}
			total.Add(int32(time.Since(s)))
			wg.Done()
		}(i)
	}
	wg.Wait()
	d = time.Since(start)
	fmt.Printf("Async: %d requests made in: %s\n", requests, d)
	fmt.Printf("Avg. Req. Duration: %s\n", time.Duration(total.Load()/int32(requests)))
}
