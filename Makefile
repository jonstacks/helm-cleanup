VERSION := 0f91717e54ebffb4f499f0004fdfe87c36fb2962

binaries:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main-linux-amd64-$(VERSION) 
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o main-linux-arm64-$(VERSION)

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

integration-test:
	go test -v -race -tags=integration .

clean:
	rm -rf ./bin
	rm -r coverage.out
