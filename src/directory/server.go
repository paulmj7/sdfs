package main

import (
	"fmt"
	"log"
	"net"
	"sdfs/services"
	"sdfs/services/pb"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	s := services.DirectoryServer{}
	gRPCServer := grpc.NewServer()
	pb.RegisterDirectoryServiceServer(gRPCServer, &s)
	fmt.Println("Listening on 9000")
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("serve error: ", err)
	}
}
