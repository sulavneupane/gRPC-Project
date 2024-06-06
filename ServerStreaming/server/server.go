package main

import (
	"fmt"
	proto "gRPC_Project/ServerStreaming/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type server struct {
	proto.UnimplementedExampleServer
}

func main() {
	listener, tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil {
		panic(tcpErr)
	}
	srv := grpc.NewServer()
	proto.RegisterExampleServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) ServerReply(req *proto.HelloRequest, stream proto.Example_ServerReplyServer) error {
	fmt.Println(req.SomeString)
	time.Sleep(5 * time.Second)
	serverReplies := make([]*proto.HelloResponse, 0)
	for i := 1; i <= 5; i++ {
		serverReplies = append(serverReplies, &proto.HelloResponse{Reply: fmt.Sprintf("Response %d", i)})
	}
	for _, serverReply := range serverReplies {
		err := stream.Send(serverReply)
		if err != nil {
			return err
		}
	}
	return nil
}
