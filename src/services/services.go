package services

import (
	"io"
	"log"
	"sdfs/directory/lib"
	"sdfs/services/pb"
	"sdfs/storage/storage_util"

	"golang.org/x/net/context"
)

// DirectoryServer struct
type DirectoryServer struct {
	pb.UnimplementedDirectoryServiceServer
}

// StorageServer struct
type StorageServer struct {
	pb.UnimplementedStorageServiceServer
}

// Register a storage server to the directory server
func (s *DirectoryServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Println("Recieve message body from client: ", in.Url)
	lib.Directory.Add(in.Url)
	lib.Directory.PrintDirectory()
	return &pb.RegisterResponse{Status: "Hello From the Server!"}, nil
}

// Create a file to the directory server
func (s *DirectoryServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Println("Receive message body from client: ", in.Name, in.Size)
	locations, err := lib.Create(in.Name, in.Size)
	if err != nil {
		log.Fatal("error creating locations: ", err)
	}
	return &pb.CreateResponse{Locations: locations}, nil
}

// Lookup chunk locations of a file
func (s *DirectoryServer) Lookup(ctx context.Context, in *pb.LookupRequest) (*pb.LookupResponse, error) {
	log.Println("Receive message body from client: ", in.Name)
	chunks, err := lib.Search(in.Name)
	if err != nil {
		log.Fatal("error searching file: ", err)
	}
	return &pb.LookupResponse{ReadChunks: chunks}, nil
}

// Ls lists all files in directory
func (s *DirectoryServer) Ls(ctx context.Context, in *pb.LsRequest) (*pb.LsResponse, error) {
	log.Println("Receive message body from client")
	names, err := lib.Ls()
	if err != nil {
		log.Fatal("error listing files: ", err)
	}
	return &pb.LsResponse{Names: names}, nil
}

// Read chunks of a file
func (s *StorageServer) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	log.Println("Hello from the client: ", in.Name)
	chunk, err := storage_util.Read(in.Name)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	return &pb.ReadResponse{Data: chunk}, nil
}

// Write chunks to a file
func (s *StorageServer) Write(stream pb.StorageService_WriteServer) error {
	log.Println("Hello From the client")
	for {
		in, err := stream.Recv()
		log.Println("Hit")
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		storage_util.Write(in.Name, in.Data)
		stream.Send(&pb.WriteResponse{Status: "Successfully wrote chunk"})
	}
	return nil
}
