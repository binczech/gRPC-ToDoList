package main

import (
	"context"
	"log"
	"os"
	"time"
	pb "todolist"

	"google.golang.org/grpc"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "1234"
	}
	address, exists := os.LookupEnv("ADDR")
	if !exists {
		address = "localhost"
	}
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure(), grpc.WithBlock())
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
	s, err := c.ReadToDo(ctx, &pb.RequestReadMessage{Id: ToDoID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("List of ToDos after Read: %s", s.GetText())
	r, err = c.DeleteToDo(ctx, &pb.DeleteToDoMessage{Id: ToDoID})
	log.Printf("List of ToDos after Delete: %s", r.GetToDosList())
}
