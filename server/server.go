package main

import (
	"io"
	"log"
	"net"

	chit "github.com/louiselykke/ChittyChat/proto"
	"google.golang.org/grpc"
)

type Server struct {
	chit.UnimplementedChatServer
	channel map[string][]chan *chit.Message
}

func (s *Server) Brodcast(msgStream chit.Chat_BrodcastServer) error {
	log.Println("Brodcast call")

	for {
		msg, err := msgStream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
		log.Printf("%s said: %s", msg.User.Name, msg.Message)
	}
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()

	chit.RegisterChatServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *Server) SendMessage(msgStream chit.Chat_BrodcastServer) error {
	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	ack := chit.Response{}
	msgStream.SendAndClose(&ack)

	go func() {
		streams := s.channel[msg.User.Name]
		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}
