package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"

	chit "github.com/louiselykke/ChittyChat/proto"
	"google.golang.org/grpc"
)

var userName = flag.String("user", "Anon", "Username for chatting")
var lamport int64 = 0

type User struct {
	id   string
	name string
}

func main() {

	flag.Parse()

	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	var thisUser = User{
		id:   "1",
		name: *userName,
	}
	//  Create new Client from generated gRPC code from proto
	c := chit.NewChatClient(conn)
	log.Printf("this is no %s", thisUser.name)
	ctx := context.Background()
	sendMessage(ctx, c, "i just joined")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, c, scanner.Text())
	}

}

func sendMessage(ctx context.Context, client chit.ChatClient, message string) {
	stream, err := client.Brodcast(ctx)
	if err != nil {
		log.Printf("Cannot send message: error: %v", err)
	}
	msg := chit.Message{
		User: &chit.User{
			Id:   "1",
			Name: *userName},
		Message: message,
		Lamport: lamport,
	}
	stream.Send(&msg)
}
