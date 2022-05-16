package storage_util

import (
	"log"
	"os"
)

func Write(chunkName string, data []byte) {
	err := os.WriteFile(chunkName, data, 0644)
	if err != nil {
		log.Fatal("error writing chunk: ", err)
	}
}

func Read(chunkName string) ([]byte, error) {
	log.Println("Chunk name: ", chunkName)
	data, err := os.ReadFile(chunkName)
	if err != nil {
		log.Fatal("error reading chunk: ", err)
	}

	log.Println("Chunk length: ", len(data))
	return data, nil
}
