build_name := asteroids

build:
	go build -o $(build_name) main.go

run:
	go run main.go

test:
	go test ./...

build-wasm:
	env GOOS=js GOARCH=wasm go build -o $(build_name).wasm main.go

clean:
	rm -f $(build_name)
	rm -f $(build_name).wasm