package main

import (
	"fmt"
	"log"
	"net"
	"sdfs/services"
	"sdfs/services/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	dir := pb.NewDirectoryServiceClient(conn)

	response, err := dir.Register(context.Background(), &pb.RegisterRequest{Url: "http://127.0.0.1:9001"})
	if err != nil {
		log.Fatal("Error when calling Register: ", err)
	}

	log.Println("Response from server: ", response.Status)

	conn.Close()

	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	s := services.StorageServer{}
	gRPCServer := grpc.NewServer(grpc.MaxRecvMsgSize(64000017), grpc.MaxSendMsgSize(64000005))
	pb.RegisterStorageServiceServer(gRPCServer, &s)
	fmt.Println("Listening on 9001")
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("serve error: ", err)
	}
}
