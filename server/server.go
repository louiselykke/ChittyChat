package main

import (
	"io"
	"log"
	"net"
	"sync"

	chit "github.com/louiselykke/ChittyChat/proto"
	"google.golang.org/grpc"
)

var lamport int = 0

type Server struct {
	chit.UnimplementedChatServer
	clients map[string]chit.Chat_BroadcastServer

	mu sync.RWMutex
}

func (s *Server) Broadcast(msgStream chit.Chat_BroadcastServer) error {
	var thisUser string
	for {
		msg, err := msgStream.Recv()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		thisUser = msg.User.Name
		defer s.removeClient(thisUser)
		s.addClient(thisUser, msgStream)

		s.updateLamportTime(int(msg.Lamport))
		log.Printf("broadcast: %s at Lamport time %d", msg.Message, lamport)
		for _, client := range s.getClients() {
			if client == msgStream { // The client sending the messages does not recieve a response through the stream.
				continue
			}
			msg.Lamport = int64(s.updateLamportTime(int(msg.Lamport)))
			if err := client.Send(msg); err != nil {
				log.Printf("broadcast err: %v", err)
			}
		}
	}
}

func (s *Server) addClient(userId string, srv chit.Chat_BroadcastServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[userId] = srv
	log.Printf("%s joined the chat! Treat them well", userId)
}

func (s *Server) removeClient(userId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, userId)
	log.Printf("%s left the chat", userId)
}

func (s *Server) getClients() []chit.Chat_BroadcastServer {
	var cs []chit.Chat_BroadcastServer

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		cs = append(cs, c)
	}
	return cs
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()

	chit.RegisterChatServer(grpcServer, &Server{clients: make(map[string]chit.Chat_BroadcastServer),
		mu: sync.RWMutex{}})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *Server) updateLamportTime(msgTime int) int {
	if lamport > msgTime {
		lamport = lamport + 1
	} else {
		lamport = msgTime + 1
	}
	return lamport
}
