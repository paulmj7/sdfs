package services

import (
	"log"
	"sdfs/directory/lib"
	"sdfs/storage/storage_util"

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
	lib.Directory.Add(in.Url)
	lib.Directory.PrintDirectory()
	return &RegisterResponse{Status: "Hello From the Server!"}, nil
}

func (s *StorageServer) Read(ctx context.Context, in *ReadRequest) (*ReadResponse, error) {
	log.Println("Hello from the client: ", in.Name)
	return &ReadResponse{Data: nil}, nil
}

func (s *StorageServer) Write(ctx context.Context, in *WriteRequest) (*WriteResponse, error) {
	log.Println("Hello From the client: ", in.Name)
	storage_util.Write(in.Name, in.Data)
	return &WriteResponse{Status: "Hello From the Server"}, nil
}
