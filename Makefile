all:
	protoc --go_out=plugins=grpc:. todolist.proto
run:
	go run server/*.go