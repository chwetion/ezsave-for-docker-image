linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/ezsave github.com/chwetion/ezsave-for-docker-image

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/macos/ezsave github.com/chwetion/ezsave-for-docker-image