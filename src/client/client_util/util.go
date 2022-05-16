package client_util

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"sdfs/services/proto"
	"strconv"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

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

func Read(fileName string) {
	chunks := lookup(fileName)
	readChunks(chunks)
}

func lookup(fileName string) []*proto.ReadChunk {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	directory := proto.NewDirectoryServiceClient(conn)
	res, err := directory.Lookup(context.Background(), &proto.LookupRequest{Name: fileName})
	if err != nil {
		log.Fatal("error searching for file: ", err)
	}
	return res.ReadChunks
}

func readChunks(chunks []*proto.ReadChunk) {
	f, err := os.Create("tempfile.mkv")
	if err != nil {
		log.Fatal("error creating temp file")
	}
	defer f.Close()

	for _, chunk := range chunks {
		data := readChunk(chunk.Name, chunk.Location)
		log.Println(len(data))
		n, err := f.Write(data)
		if err != nil {
			log.Fatal("error writing file: ", err, n)
		}
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

	storage := proto.NewStorageServiceClient(conn)

	res, err := storage.Read(context.Background(), &proto.ReadRequest{Name: chunkName})
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

	directory := proto.NewDirectoryServiceClient(conn)

	res, err := directory.Create(context.Background(), &proto.CreateRequest{Name: fileName, Size: fileSize})
	if err != nil {
		log.Fatal("error creating file: ", err)
	}
	return res.Locations
}

func write(fileName string, locations []string) {
	f, err := os.Open(fileName)
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
		wg.Add(1)
		go writeChunk(chunkName, buf[:n], locations[i])
		wg.Wait()
	}
}

func writeChunk(chunkName string, data []byte, location string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(location, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	storage := proto.NewStorageServiceClient(conn)

	defer wg.Done()
	response, err := storage.Write(context.Background(), &proto.WriteRequest{Name: chunkName, Data: data})
	if err != nil {
		log.Fatal("Error when calling Write: ", err)
	}

	log.Println("Response from server: ", response.Status)
}
