package main

import (
	"log"
	"sdfs/services"

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

	dir := services.NewDirectoryServiceClient(conn)

	response, err := dir.Register(context.Background(), &services.RegisterRequest{Url: "Hello From Client!"})
	if err != nil {
		log.Fatal("Error when calling Register: ", err)
	}

	log.Println("Response from server: ", response.Status)
}
