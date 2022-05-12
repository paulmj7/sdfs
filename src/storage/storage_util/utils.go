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
