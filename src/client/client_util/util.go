package client_util

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"sdfs/services"
	"strconv"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

func Partition(fileName string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: ", err)
	}
	defer conn.Close()

	storage := services.NewStorageServiceClient(conn)

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
		go writeChunk(chunkName, buf[:n], storage)
		wg.Wait()
	}
}

func writeChunk(chunkName string, data []byte, storage services.StorageServiceClient) {
	defer wg.Done()
	response, err := storage.Write(context.Background(), &services.WriteRequest{Name: chunkName, Data: data})
	if err != nil {
		log.Fatal("Error when calling Write: ", err)
	}

	log.Println("Response from server: ", response.Status)
}
