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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	text := "name"
	r, err := c.ListToDos(ctx, &pb.RequestListMessage{Text: text})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List ToDos: %s", r.GetToDosList())
	r, err = c.AddToDo(ctx, &pb.AddToDoMessage{Text: "New ToDo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List ToDos after Add: %s", r.GetToDosList())
	ToDosList := r.GetToDosList()
	var ToDoID string
	for _, ToDo := range ToDosList {
		if ToDo.Text == "New ToDo" {
			ToDoID = ToDo.Id
		}
	}
	r, err = c.UpdateToDo(ctx, &pb.UpdateToDoMessage{Id: ToDoID, Text: "Updated ToDo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List ToDos after Update: %s", r.GetToDosList())
	r, err = c.DeleteToDo(ctx, &pb.DeleteToDoMessage{Id: ToDoID})
	log.Printf("List of ToDos after Delete: %s", r.GetToDosList())
}
