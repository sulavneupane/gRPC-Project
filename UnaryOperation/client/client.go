package main

import (
	"context"
	proto "gRPC_Project/UnaryOperation/protoc"
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
	r.GET("/send-message-to-server/:message", clientConnectionServer)
	_ = r.Run(":8000")
}

func clientConnectionServer(c *gin.Context) {
	message := c.Param("message")

	req := &proto.HelloRequest{SomeString: message}

	_, _ = client.ServerReply(context.TODO(), req)
	c.JSON(http.StatusOK, gin.H{
		"message": "message sent successfully to server " + message,
	})
}
