package main

import (
	"errors"
	"fmt"
	proto "gRPC_Project/BiDirectionalStreaming/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"strconv"
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

func (s *server) ServerReply(stream proto.Example_ServerReplyServer) error {
	for i := 0; i < 5; i++ {
		err := stream.Send(&proto.HelloResponse{Reply: "message " + strconv.Itoa(i) + " from server"})
		if err != nil {
			return errors.New("unable to send data from server")
		}
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(req.SomeString)
	}
	return nil
}
