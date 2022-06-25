package client_util

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"sdfs/services/pb"
	"strconv"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

type storageClientWrapper struct {
	Client pb.StorageService_WriteClient
}

// Ls lists all files in directory
func Ls() {

}

// Create file in directory
func Create(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal("error getting file info: ", err)
	}
	locations := create(fileName, uint64(fi.Size()))
	if len(locations) == 0 {
		log.Fatal("error getting locations")
	}
	write(fileName, locations)
}

// Read file from directory
func Read(fileName string) {
	chunks := lookup(fileName)
	readChunks(chunks)
}

func lookup(fileName string) []*pb.ReadChunk {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	directory := pb.NewDirectoryServiceClient(conn)
	res, err := directory.Lookup(context.Background(), &pb.LookupRequest{Name: fileName})
	if err != nil {
		log.Fatal("error searching for file: ", err)
	}
	return res.ReadChunks
}

func readChunks(chunks []*pb.ReadChunk) {
	for _, chunk := range chunks {
		data := readChunk(chunk.Name, chunk.Location)
		os.Stdout.Write(data)
	}
}

// pretty sure the issue is that the ReadRequest is a weird memory thing, need to refactor back into services package

func readChunk(chunkName, location string) []byte {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(location, grpc.WithInsecure(), grpc.WithMaxMsgSize(64000005))
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	storage := pb.NewStorageServiceClient(conn)

	res, err := storage.Read(context.Background(), &pb.ReadRequest{Name: chunkName})
	if err != nil {
		log.Fatal("error reading chunk client", err)
	}

	return res.Data
}

func create(fileName string, fileSize uint64) []string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	directory := pb.NewDirectoryServiceClient(conn)

	res, err := directory.Create(context.Background(), &pb.CreateRequest{Name: fileName, Size: fileSize})
	if err != nil {
		log.Fatal("error creating file: ", err)
	}
	return res.Locations
}

func write(fileName string, locations []string) {
	f, err := os.Open(fileName)
	streams := make(map[string]*storageClientWrapper)
	if err != nil {
		log.Fatal("error opening file: ", err)
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 64000000)

	buf := make([]byte, 64000000)
	h := sha1.New()
	for i := 0; ; i++ {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error buffering file: ", err)
		}
		io.WriteString(h, fileName+strconv.Itoa(i))
		chunkName := hex.EncodeToString(h.Sum(nil))[:10]
		stream, exists := streams[locations[i]]
		var conn *grpc.ClientConn
		if !exists {
			stream, conn = setupStorageClient(locations[i])
			// is this secure?
			defer conn.Close()
			streams[locations[i]] = stream
		}
		wg.Add(1)
		go writeChunk(chunkName, buf[:n], stream)
		wg.Wait()
	}
}

func setupStorageClient(location string) (*storageClientWrapper, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(location, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Fatal("did not connect: ", err)
	}

	client := pb.NewStorageServiceClient(conn)
	stream, err := client.Write(context.Background())
	wrapper := storageClientWrapper{Client: stream}
	return &wrapper, conn
}

func writeChunk(chunkName string, data []byte, stream *storageClientWrapper) {
	stream.Client.Send(&pb.WriteRequest{Name: chunkName, Data: data})
	go func() {
		response, err := stream.Client.Recv()
		defer wg.Done()
		if err != nil {
			log.Fatal("Error writing chunk: ", err)
		}
		log.Println(response.Status)
	}()
}
