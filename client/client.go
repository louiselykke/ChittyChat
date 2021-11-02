package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	chit "github.com/louiselykke/ChittyChat/proto"
	"google.golang.org/grpc"
)

var userName = flag.String("user", "Anon", "Username for chatting")
var lamport int = 0

func main() {

	flag.Parse() // adds the value of the flag parameters default or userspecific

	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	x := rand.NewSource(time.Now().UnixNano())
	y := rand.New(x)

	user := chit.User{
		Id:   strconv.Itoa(y.Intn(1000)),
		Name: *userName,
	}

	// When this function returns, call con.Close
	// use defer to do so
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := chit.NewChatClient(conn)
	welcome()
	ctx := context.Background()

	// creating a bidirectional stream / but this means that we cannot join and cut the connection with a function call :(

	stream, err := c.Broadcast(ctx) // this establishes the connection to the server.
	if err != nil {
		log.Fatal(err)
	}
	if err := stream.SendMsg(&chit.Message{ // send the initial message.
		User:    &user,
		Message: "Attempting to register " + user.Name + " on the server",
		Lamport: int64(lamport + 1),
	}); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() { /// for each text entered in the terminal a message will be send to the server through
			msg := scanner.Text()
			lamport = lamport + 1
			if err := stream.SendMsg(&chit.Message{
				User:    &user,
				Message: msg,
				Lamport: int64(lamport),
			}); err != nil {
				log.Fatal(err)
			}
			log.Printf("sent: %s, at lamport time %d", msg, lamport)
		}
	}()
	for {
		resp, err := stream.Recv() //
		if err != nil {
			log.Fatal(err)
		}

		updateLamportTime(int(resp.Lamport))

		// add ifstatement to no log message if this client just send that message...

		log.Printf("recv from %s: %s, at Lamport time %d", resp.User.Name, resp.Message, resp.Lamport)
	}
}

func welcome() {
	log.Println(`
	__________________
	< CHITTY CHAT CHAT >
	 ------------------`)
}

func updateLamportTime(msgTime int) int {
	if lamport > msgTime {
		lamport = lamport + 1
	} else {
		lamport = msgTime + 1
	}
	return lamport
}

//#### old methods not relevant
func publish(ctx context.Context, client chit.Chat_BroadcastClient, message string) {

	msg := chit.Message{
		User: &chit.User{
			Id:   "1",
			Name: *userName},
		Message: message,
		Lamport: int64(lamport),
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
