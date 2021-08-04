# Cross-platform Example

This is a basic example of constructing a controllerbus, running a controller,
and executing a directive against that controller. It runs in all platforms that
Go or TinyGo support, including the web browser with WebAssembly or GopherJS.

## Running on Desktop

To run on a host with the Go tool:

```sh
go run ./
```

## Running in Browser with WebAssembly

Build the WebAssembly payload and run the server:

```sh
./wasm.bash
```

Browse to localhost:5000, open the console to view the log messages, and click
"Run" to start the program.

## Running in Browser with GopherJS

Install the latest gopherjs release:

```sh
# run in ~
go install github.com/gopherjs/gopherjs@master
```

Build the GopherJS payload and run the server:

```sh
./gopherjs.bash
```

Browse to http://localhost:8080/github.com/paralin/go-indexeddb/example/
