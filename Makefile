BUILD_PATH := $(realpath ./build/server)

chmod-deploy-scripts:
	chmod +x ./deploy/start.sh

executable-path:
	@echo $(BUILD_PATH)

dev:
	go run ./cmd/server/*.go

build-server:
	mkdir -p build
	CGO_ENABLED=0 go build -o ./build/server -v ./cmd/server/*.go

chmod-server:
	chmod +x ./build/server

build: build-server chmod-server
