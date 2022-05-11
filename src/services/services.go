package services

import (
	"log"

	"golang.org/x/net/context"
)

type DirectoryServer struct {
	UnimplementedDirectoryServiceServer
}

type StorageServer struct {
	UnimplementedStorageServiceServer
}

func (s *DirectoryServer) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	log.Println("Recieve message body from client: ", in.Url)
	return &RegisterResponse{Status: "Hello From the Server!"}, nil
}

func (s *StorageServer) Read(ctx context.Context, in *ReadRequest) (*ReadResponse, error) {
	log.Println("Hello from the client: ", in.Name)
	return &ReadResponse{Data: nil}, nil
}

func (s *StorageServer) Write(ctx context.Context, in *WriteRequest) (*WriteResponse, error) {
	log.Println("Hello From the client: ", in.Name)
	return &WriteResponse{Status: "Hello From the Server"}, nil
}
