package main

import (
	"context"
	"log"
	"time"
	pb "todolist"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewToDoListManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	text := "name"
	r, err := c.ListToDos(ctx, &pb.RequestListMessage{Text: text})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List of ToDos: %s", r.GetToDosList())
	r, err = c.AddToDo(ctx, &pb.AddToDoMessage{Text: "New ToDo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List of ToDos: %s", r.GetToDosList())
}
