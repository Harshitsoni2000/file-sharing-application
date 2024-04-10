package main

import (
	pb "github.com/Harshitsoni2000/File_Sharing_Application/server/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	chunkSize = 3670016
)

type FileServer struct {
	pb.FileServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error while starting tcp server :: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, &FileServer{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
