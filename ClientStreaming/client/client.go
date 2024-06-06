package main

import (
	"context"
	"fmt"
	proto "gRPC_Project/ClientStreaming/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

var client proto.ExampleClient

func main() {
	// Connection to internal grpc server
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewExampleClient(conn)

	// Implement REST API
	r := gin.Default()
	r.GET("/send", clientConnectionServer)
	_ = r.Run(":8000")
}

func clientConnectionServer(c *gin.Context) {
	reqs := make([]*proto.HelloRequest, 0)
	for i := 1; i <= 6; i++ {
		reqs = append(reqs, &proto.HelloRequest{SomeString: fmt.Sprintf("Request %d", i)})
	}

	stream, err := client.ServerReply(context.TODO())
	if err != nil {
		fmt.Println("something went wrong", err)
	}
	for _, req := range reqs {
		err = stream.Send(req)
		if err != nil {
			fmt.Println("request not fulfilled", err)
			return
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("there is some error occurred", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message_count": response,
	})
}
