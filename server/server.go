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
var clients map[string]chit.Chat_BroadcastServer

type Server struct {
	chit.UnimplementedChatServer

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
		s.updateLamportTime(int(msg.Lamport))
		thisUser = msg.User.Name
		log.Printf("broadcast: %s at Lamport time %d", msg.Message, lamport)
		if clients[thisUser] == nil {
			s.addClient(thisUser, msgStream)
			defer s.removeClient(thisUser)
		}
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

func (s *Server) addClient(userName string, srv chit.Chat_BroadcastServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateLamportTime(lamport)
	clients[userName] = srv
	log.Printf("%s joined the chat, at Lamport time %d! Treat them well", userName, lamport)

}

func (s *Server) removeClient(userId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := clients[userId]; ok {
		delete(clients, userId)
		s.updateLamportTime(lamport)
		log.Printf("%s left the chat at Lamport time %d", userId, lamport)

	}
}

func (s *Server) getClients() []chit.Chat_BroadcastServer {
	var cs []chit.Chat_BroadcastServer

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range clients {
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

	clients = make(map[string]chit.Chat_BroadcastServer)

	chit.RegisterChatServer(grpcServer, &Server{
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
