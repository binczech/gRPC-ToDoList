all:
	protoc --go_out=plugins=grpc:. todolist.proto