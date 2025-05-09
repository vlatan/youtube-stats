CLI_NAME=cli
WEB_NAME=ui

build:
	go build -o ${CLI_NAME} ./cmd/cli
	go build -o ${WEB_NAME} ./cmd/web

run: build
	./${WEB_NAME}

clean:
	go clean
	rm ${CLI_NAME}
	rm ${WEB_NAME}