# Compute-Go-VM

Quick example of instantiating a gRPC Go service in a custom container at runtime

## Prerequisites

* Docker is installed and running

## Instructions

1. Run `make build-worker`
2. Run `make run`

## Notes

* Need to run `go run cmd/main/main.go` as `root` or else it won't have the necessary permissions to bind to the unix socket made by the VM

## TODO

* Clean up unix socket after container termination
