package main

import (
    "log"
    "net"
    "os"

    "google.golang.org/grpc"
    pb "homework/pb"
    "homework/server"
)

const uploadDirectory = "./uploads/"

func main() {
    // Create the upload directory if it doesn't exist
    if err := os.MkdirAll(uploadDirectory, os.ModePerm); err != nil {
        log.Fatalf("Failed to create upload directory: %v", err)
    }

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterFileUploadServer(s, &server.FileUploadServer{})
    log.Printf("Server started on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
