package main

import (
	"context"
	"fmt"
	proto "gRPC_Project/ServerStreaming/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"time"
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
	stream, err := client.ServerReply(context.TODO(), &proto.HelloRequest{
		SomeString: "Client Request",
	})
	if err != nil {
		fmt.Println("something went wrong", err)
		return
	}

	count := 0
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("server response: ", response)
		time.Sleep(1 * time.Second)
		count++
	}
	c.JSON(http.StatusOK, gin.H{
		"message_count": count,
	})
}
