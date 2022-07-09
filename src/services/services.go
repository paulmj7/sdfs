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
	lib.Directory.Add(in.Url)
	lib.Directory.PrintDirectory()
	return &pb.RegisterResponse{Status: "Hello From the Server!"}, nil
}

// Create a file to the directory server
func (s *DirectoryServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	locations, err := lib.Create(in.Name, in.Size)
	if err != nil {
		log.Fatal("error creating locations: ", err)
	}
	return &pb.CreateResponse{Locations: locations}, nil
}

// Lookup chunk locations of a file
func (s *DirectoryServer) Lookup(ctx context.Context, in *pb.LookupRequest) (*pb.LookupResponse, error) {
	chunks, err := lib.Search(in.Name)
	if err != nil {
		log.Fatal("error searching file: ", err)
	}
	return &pb.LookupResponse{ReadChunks: chunks}, nil
}

// Ls lists all files in directory
func (s *DirectoryServer) Ls(ctx context.Context, in *pb.LsRequest) (*pb.LsResponse, error) {
	names, err := lib.Ls()
	if err != nil {
		log.Fatal("error listing files: ", err)
	}
	return &pb.LsResponse{Names: names}, nil
}

// Rm remove file from directory
func (s *DirectoryServer) Rm(ctx context.Context, in *pb.RmRequest) (*pb.RmResponse, error) {
	err := lib.Rm(in.Name)
	if err != nil {
		log.Fatal("error removing file: ", err)
		return &pb.RmResponse{Status: "Failure"}, err
	}
	return &pb.RmResponse{Status: "Success"}, nil
}

// Read chunks of a file
func (s *StorageServer) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadResponse, error) {
	chunk, err := storage_util.Read(in.Name)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	return &pb.ReadResponse{Data: chunk}, nil
}

// Write chunks to a file
func (s *StorageServer) Write(stream pb.StorageService_WriteServer) error {
	for {
		in, err := stream.Recv()
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

// Delete chunks from storage server
func (s *StorageServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := storage_util.Delete(in.Names)
	if err != nil {
		log.Fatal("error deleting chunks: ", err)
		return &pb.DeleteResponse{Status: "Failure"}, err
	}
	return &pb.DeleteResponse{Status: "Success"}, nil
}
