FILENAME=fim
VERSION=v1.1.1

LDFLAGS=-ldflags "-X main.BuildVersion=$(VERSION)"

windows:
	GOARCH=amd64 GOOS=windows go build $(LDFLAGS) -o dist/$(FILENAME).exe

linux:
	GOARCH=amd64 GOOS=linux go build $(LDFLAGS) -o dist/$(FILENAME)

wasm:
	GOARCH=wasm GOOS=js go build $(LDFLAGS) -o dist/$(FILENAME).wasm
