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

func main() {

	flag.Parse() // adds the value of the flag parameters default or userspecific

	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// When this function returns, call con.Close
	// use defer to do so
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := chit.NewChatClient(conn)
	welcome()
	ctx := context.Background()

	// creating a bidirectional stream / but this means that we cannot join and cut the connection with a function call :(
	stream, err := c.Broadcast(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := stream.SendMsg(&chit.Message{ //
		User: &chit.User{
			Id:   "1",
			Name: *userName},
		Message: "",
		Lamport: lamport,
	}); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() {
			msg := scanner.Text()
			if err := stream.SendMsg(&chit.Message{
				User: &chit.User{
					Id:   "1",
					Name: *userName},
				Message: msg,
				Lamport: lamport,
			}); err != nil {
				log.Fatal(err)
			}
			log.Printf("sent: %s", msg)
		}
	}()
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		// add ifstatement to no log message if this client just send that message...

		log.Printf("recv from %s: %s", resp.User.Name, resp.Message)
	}
}

func welcome() {
	log.Println(`
	__________________
	< CHITTY CHAT CHAT >
	 ------------------`)
}

//#### old methods not relevant
func publish(ctx context.Context, client chit.Chat_BroadcastClient, message string) {

	msg := chit.Message{
		User: &chit.User{
			Id:   "1",
			Name: *userName},
		Message: message,
		Lamport: lamport,
	}
	client.Send(&msg)
}

func recieve(ctx context.Context, client chit.Chat_BroadcastClient) {
	for {
		msg, err := client.Recv()
		if err != nil {
			log.Printf("Oops nothing recived: %v", err)
		}
		log.Printf("%s says: %s -- at Lamport time %b", msg.GetUser(), msg.GetMessage(), msg.Lamport)
	}
}
