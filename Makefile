.PHONY: protos
protos:
	mkdir -p internal/generated; protoc --go_out=internal/generated --go_opt=paths=source_relative \
    --go-grpc_out=internal/generated --go-grpc_opt=paths=source_relative \
    protos/**/*.proto


.PHONY: build-worker
build-worker: protos
	docker build -t compute:worker-go -f Dockerfile.worker .

.PHONY: build-js-worker
build-js-worker:
	docker build -t compute:worker-js -f js_worker/Dockerfile .

.PHONY: build-main
build-main: protos
	docker build -t  compute:orchestration -f Dockerfile .

.PHONY: build
build: build-worker build-js-worker build-main

.PHONY: run
run:
	docker run -v /var/run/docker.sock:/var/run/docker.sock -v /tmp/compute:/tmp/compute compute:orchestration

ssh:
	docker run -v /var/run/docker.sock:/var/run/docker.sock -v /tmp/compute:/tmp/compute -it --entrypoint bash compute:orchestration

