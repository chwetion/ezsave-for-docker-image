VERSION=1.0.0

default: linux macos

linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/ezsave-linux-amd64-v$(VERSION) github.com/chwetion/ezsave-for-docker-image

macos: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/macos/ezsave-darwin-amd64-v$(VERSION) github.com/chwetion/ezsave-for-docker-image

clean:
	rm -rf bin/linux/*
	rm -rf bin/macos/*