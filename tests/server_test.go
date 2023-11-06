package test

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	pb "homework/pb"
	"homework/server"

	"google.golang.org/grpc"
)

const testUploadDirectory = "./test_uploads/"

func setupTest(t *testing.T) {
	if err := os.MkdirAll("test_uploads", os.ModePerm); err != nil {
		t.Fatalf("Failed to create test uploads directory: %v", err)
	}
}

func TestUploadFile(t *testing.T) {

	setupTest(t)

	s := grpc.NewServer()
	pb.RegisterFileUploadServer(s, &server.FileUploadServer{})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	go s.Serve(listener)

	defer s.Stop()

	client, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	clientService := pb.NewFileUploadClient(client)

	files, err := os.ReadDir(testUploadDirectory)
	if err != nil {
		t.Fatalf("Failed to read the directory: %v", err)
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		uploadStream, err := clientService.UploadFile(context.Background())
		if err != nil {
			t.Fatalf("Failed to create upload stream: %v", err)
		}

		filePath := filepath.Join(testUploadDirectory, file.Name())
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read the test file: %v", err)
		}

		chunkSize := 64

		for i := 0; i < len(fileData); i += chunkSize {
			end := i + chunkSize
			if end > len(fileData) {
				end = len(fileData)
			}

			chunk := &pb.FileChunk{Data: fileData[i:end], FileName: file.Name()}
			if err := uploadStream.Send(chunk); err != nil {
				t.Fatalf("Error sending chunk: %v", err)
			}
		}

		if err := uploadStream.CloseSend(); err != nil {
			t.Fatalf("Error closing the stream: %v", err)
		}

		response, err := uploadStream.CloseAndRecv()
		if err != nil {
			t.Fatalf("Error receiving response: %v", err)
		}

		if !response.Success {
			t.Fatalf("Upload failed for file: %s", file.Name())
		}

		log.Printf("File uploaded successfully: %s", file.Name())
	}
}
