**About**

Simple API for managing To Do List. Consists of a server and a client. Server implements functions for listing, creating, updating and deleting To Dos. Client implements example of communication with the server.

**Prerequisities**
* [Golang](https://golang.org/doc/install)
* [Protocol](https://developers.google.com/protocol-buffers)
* `export GO111MODULE=on  # Enable module mode`
* `export PATH="$PATH:$(go env GOPATH)/bin"`
* `export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/server/binczech-test-a273644ddbb5.json"`

**How to run**
1. `make`
2. In first terminal run `go run server/*.go`
3. In second terminal run `go run client/*.go`