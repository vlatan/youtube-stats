CLI_NAME=cli
WEB_NAME=ui

build:
	go build -o ${CLI_NAME} cmd/cli/cli.go
	go build -o ${WEB_NAME} cmd/web/web.go

run: build
	./${WEB_NAME}

clean:
	go clean
	rm ${CLI_NAME}
	rm ${WEB_NAME}