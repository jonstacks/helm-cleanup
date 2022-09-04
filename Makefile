VERSION := 48044648aaf91dd892f355543ab16f931af594c4

binaries:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main-linux-amd64-$(VERSION) 
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o main-linux-arm64-$(VERSION)

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-integration:
	go test -v -race -timeout=15m -tags=integration .

clean:
	rm -rf ./bin
	rm -r coverage.out
