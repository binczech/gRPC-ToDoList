package main

import (
	"context"
	"log"
	"net"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"

	pb "todolist"

	firebase "firebase.google.com/go"
)

type server struct {
	Client *firestore.Client
	pb.UnimplementedToDoListManagerServer
}

func connect() (client *firestore.Client) {
	ctx := context.Background()
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

func (s *server) getList(ctx context.Context) (ToDosList []*pb.ToDoMessage) {
	iter := s.Client.Collection("todolist").Documents(ctx)
	for {
		var fetchedData *pb.ToDoMessage
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
		}
		if err := doc.DataTo(&fetchedData); err != nil {
			log.Printf("Failed to iterate: %v", err)
		}
		fetchedData.Id = doc.Ref.ID
		ToDosList = append(ToDosList, fetchedData)
	}
	return ToDosList
}

func (s *server) setList(ctx context.Context, in *pb.AddToDoMessage) {
	ref := s.Client.Collection("todolist").NewDoc()
	_, err := ref.Set(ctx, in)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
}

func (s *server) deleteFromList(ctx context.Context, id string) {
	ref := s.Client.Collection("todolist")
	toDelete := ref.Doc(id)
	_, err := toDelete.Delete(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
}

func (s *server) getOneDoc(id string) (ToDo *firestore.DocumentRef) {
	ref := s.Client.Collection("todolist")
	ToDo = ref.Doc(id)
	return ToDo
}

func (s *server) updateInList(ctx context.Context, in *pb.UpdateToDoMessage) {
	idToUpdate := in.GetId()
	text := in.GetText()
	toUpdate := s.getOneDoc(idToUpdate)
	_, err := toUpdate.Set(ctx, &pb.UpdateToDoMessage{
		Text: text,
	})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
}

func (s *server) ReadToDo(ctx context.Context, in *pb.RequestReadMessage) (*pb.ToDoMessage, error) {
	idToRead := in.GetId()
	toRead := s.getOneDoc(idToRead)
	docsnap, err := toRead.Get(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	dataMap := docsnap.Data()
	return &pb.ToDoMessage{Id: dataMap["Id"].(string), Text: dataMap["Text"].(string)}, nil
}

func (s *server) ListToDos(ctx context.Context, in *pb.RequestListMessage) (*pb.ListToDosMessage, error) {
	// Use the application default credentials
	ToDosList := s.getList(ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func (s *server) AddToDo(ctx context.Context, in *pb.AddToDoMessage) (*pb.ListToDosMessage, error) {
	s.setList(ctx, in)
	ToDosList := s.getList(ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func (s *server) DeleteToDo(ctx context.Context, in *pb.DeleteToDoMessage) (*pb.ListToDosMessage, error) {
	idToDelete := in.GetId()
	s.deleteFromList(ctx, idToDelete)
	ToDosList := s.getList(ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func (s *server) UpdateToDo(ctx context.Context, in *pb.UpdateToDoMessage) (*pb.ListToDosMessage, error) {
	s.updateInList(ctx, in)
	ToDosList := s.getList(ctx)
	return &pb.ListToDosMessage{ToDosList: ToDosList}, nil
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "1234"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := server{Client: connect()}
	defer server.Client.Close()
	pb.RegisterToDoListManagerServer(s, &server)
	log.Printf("Listening on %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
