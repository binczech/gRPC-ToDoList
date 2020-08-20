package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"

	pb "todolist"

	firebase "firebase.google.com/go"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedToDoListManagerServer
}

func Connect(ctx context.Context) (client *firestore.Client) {
	ctx = context.Background()
	conf := &firebase.Config{ProjectID: "binczech-test"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func GetList(client *firestore.Client, ctx context.Context) (ToDosList []*pb.ToDoMessage) {
	iter := client.Collection("todolist").Documents(ctx)
	for {
		var fetchedData *pb.ToDoMessage
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		if err := doc.DataTo(&fetchedData); err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fetchedData.Id = doc.Ref.ID
		fmt.Println(fetchedData)
		ToDosList = append(ToDosList, fetchedData)
	}
	fmt.Println(ToDosList)
	return ToDosList
}

func SetList(client *firestore.Client, ctx context.Context, in *pb.AddToDoMessage) {
	ref := client.Collection("todolist").NewDoc()
	_, err := ref.Set(ctx, in)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
}

func (s *server) ListToDos(ctx context.Context, in *pb.RequestListMessage) (*pb.ListToDosMessage, error) {
	// Use the application default credentials
	client := Connect(ctx)
	defer client.Close()
	ToDosList := GetList(client, ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func (s *server) AddToDo(ctx context.Context, in *pb.AddToDoMessage) (*pb.ListToDosMessage, error) {
	client := Connect(ctx)
	SetList(client, ctx, in)
	ToDosList := GetList(client, ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterToDoListManagerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
