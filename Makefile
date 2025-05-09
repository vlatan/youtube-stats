CLI_NAME=cli
WEB_NAME=app

build:
	go build -o ./bin/${CLI_NAME} ./cmd/cli
	go build -o ./bin/${WEB_NAME} ./cmd/app

run: build
	./bin/${WEB_NAME}

clean:
	go clean
	rm ./bin/${CLI_NAME}
	rm ./bin/${WEB_NAME}