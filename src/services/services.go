package services

import (
	"log"
	"sdfs/directory/lib"
	"sdfs/services/proto"
	"sdfs/storage/storage_util"

	"golang.org/x/net/context"
)

type DirectoryServer struct {
	proto.UnimplementedDirectoryServiceServer
}

type StorageServer struct {
	proto.UnimplementedStorageServiceServer
}

func (s *DirectoryServer) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	log.Println("Recieve message body from client: ", in.Url)
	lib.Directory.Add(in.Url)
	lib.Directory.PrintDirectory()
	return &proto.RegisterResponse{Status: "Hello From the Server!"}, nil
}

func (s *DirectoryServer) Create(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	log.Println("Receive message body from client: ", in.Name, in.Size)
	locations, err := lib.Create(in.Name, in.Size)
	if err != nil {
		log.Fatal("error creating locations: ", err)
	}
	return &proto.CreateResponse{Locations: locations}, nil
}

func (s *DirectoryServer) Lookup(ctx context.Context, in *proto.LookupRequest) (*proto.LookupResponse, error) {
	log.Println("Receive message body from client: ", in.Name)
	chunks, err := lib.Search(in.Name)
	if err != nil {
		log.Fatal("error searching file: ", err)
	}
	return &proto.LookupResponse{ReadChunks: chunks}, nil
}

func (s *StorageServer) Read(ctx context.Context, in *proto.ReadRequest) (*proto.ReadResponse, error) {
	log.Println("Hello from the client: ", in.Name)
	chunk, err := storage_util.Read(in.Name)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	return &proto.ReadResponse{Data: chunk}, nil
}

func (s *StorageServer) Write(ctx context.Context, in *proto.WriteRequest) (*proto.WriteResponse, error) {
	log.Println("Hello From the client: ", in.Name)
	storage_util.Write(in.Name, in.Data)
	return &proto.WriteResponse{Status: "Hello From the Server"}, nil
}
