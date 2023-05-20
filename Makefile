.PHONY: protos
protos:
	mkdir -p internal/generated; protoc --go_out=internal/generated --go_opt=paths=source_relative \
    --go-grpc_out=internal/generated --go-grpc_opt=paths=source_relative \
    protos/**/*.proto


.PHONY: build-worker
build-worker: protos
	docker build -t compute:worker01 -f Dockerfile.worker .

.PHONY: run
run:
	sudo go run cmd/main/main.go