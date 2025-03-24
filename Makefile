FILENAME=fim

windows:
	GOARCH=amd64 GOOS=windows go build -o dist/$(FILENAME).exe

linux:
	GOARCH=amd64 GOOS=linux go build -o dist/$(FILENAME)

wasm:
	GOARCH=wasm GOOS=js go build -o dist/$(FILENAME).wasm
