package main

import (
	"fmt"
	proto "gRPC_Project/BigFilesStreaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"os"
)

type server struct {
	proto.UnimplementedStreamUploadServer
}

func main() {
	listener, tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil {
		panic(tcpErr)
	}
	srv := grpc.NewServer()
	proto.RegisterStreamUploadServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s server) Upload(stream proto.StreamUpload_UploadServer) error {
	var fileBytes []byte
	var fileSize int64 = 0

	var fileName string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fileName = req.GetFilePath()
		chunks := req.GetChunks()
		fileBytes = append(fileBytes, chunks...)
		fileSize += int64(len(chunks))
	}

	fmt.Println(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err2 := f.Write(fileBytes)
	if err2 != nil {
		return err2
	}

	return stream.SendAndClose(&proto.UploadResponse{
		FileSize: fileSize,
		Message:  "File written successfully",
	})
}
