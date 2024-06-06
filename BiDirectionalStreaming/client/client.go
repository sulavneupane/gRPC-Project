package main

import (
	"context"
	"fmt"
	proto "gRPC_Project/BiDirectionalStreaming/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"strconv"
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
	stream, err := client.ServerReply(context.TODO())
	if err != nil {
		fmt.Println("something went wrong", err)
		return
	}

	send, receive := 0, 0
	for i := 0; i < 10; i++ {
		err := stream.Send(&proto.HelloRequest{SomeString: "message " + strconv.Itoa(i) + " from client"})
		if err != nil {
			fmt.Println("unable to send data", err)
			return
		}
		send++
	}
	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("Server message: ", message.Reply)
		receive++
	}
	c.JSON(http.StatusOK, gin.H{
		"message_sent":     send,
		"message_received": receive,
	})
}
