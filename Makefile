APP_BINARY=app
CLI_BINARY=cli
BUILD_FLAGS = CGO_ENABLED=0 GOOS=linux GOARCH=amd64

.PHONY: build
build:
	@echo "building all binaries..."
	@${BUILD_FLAGS} go build -o ./bin/${APP_BINARY} ./cmd/app
	@${BUILD_FLAGS} go build -o ./bin/${CLI_BINARY} ./cmd/cli

.PHONY: build-app
build-app:
	@echo "building app..."
	@${BUILD_FLAGS} go build -o ./bin/${APP_BINARY} ./cmd/app

.PHONY: build-cli
build-cli:
	@echo "building cli..."
	@${BUILD_FLAGS} go build -o ./bin/${CLI_BINARY} ./cmd/cli

.PHONY: run-app 
run-app: build-app
	@./bin/${APP_BINARY}

.PHONY: clean 
clean:
	@echo "cleaning all binaries..."
	@go clean
	@rm -f ./bin/*
	
