package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	chatpb "github.com/TarlexGit/go-grpc-chat/pb"
	// p "github.com/TarlexGit/go-grpc-chat/pb"
	// "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// type MessageServer struct {
// 	// grpc.UnsafeFooBarServiceServer
// 	p.UnsafeMessageServiceServer
// }

// var port = ":8080"

// func (MessageServer) SayIt(ctx context.Context, r *p.Request) (*p.Response, error) {
// 	fmt.Println("Request Text:", r.Text)
// 	fmt.Println("Request SubText:", r.Subtext)
// 	response := &p.Response{
// 		Text:    r.Text,
// 		Subtext: "Got it!",
// 	}
// 	return response, nil
// }

type chatServiceServer struct {
	chatpb.UnimplementedChatServiceServer
	mu      sync.Mutex
	channel map[string][]chan *chatpb.Message
}

func (s *chatServiceServer) JoinChannel(ch *chatpb.Channel, msgStream chatpb.ChatService_JoinChannelServer) error {

	msgChannel := make(chan *chatpb.Message)
	s.channel[ch.Name] = append(s.channel[ch.Name], msgChannel)

	// doing this never closes the stream
	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-msgChannel:
			fmt.Printf("GO ROUTINE (got message): %v \n", msg)
			msgStream.Send(msg)
		}
	}
}

func (s *chatServiceServer) SendMessage(msgStream chatpb.ChatService_SendMessageServer) error {
	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	ack := chatpb.MessageAck{Status: "SENT"}
	msgStream.SendAndClose(&ack)

	go func() {
		streams := s.channel[msg.Channel.Name]
		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}

func newServer() *chatServiceServer {
	s := &chatServiceServer{
		channel: make(map[string][]chan *chatpb.Message),
	}
	fmt.Println(s)
	return s
}

func main() {
	fmt.Println("--- SERVER APP ---")
	lis, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatpb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
	// server := grpc.NewServer()
	// var messageServer MessageServer

	// p.RegisterMessageServiceServer(server, messageServer)
	// listen, err := net.Listen("tcp", port)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Serving requests...")
	// server.Serve(listen)
}
