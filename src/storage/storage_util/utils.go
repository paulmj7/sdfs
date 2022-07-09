package storage_util

import (
	"log"
	"os"
)

// Write chunk to storage server
func Write(chunkName string, data []byte) {
	err := os.WriteFile(chunkName, data, 0644)
	if err != nil {
		log.Fatal("error writing chunk: ", err)
	}
}

// Read chunk from storage server
func Read(chunkName string) ([]byte, error) {
	data, err := os.ReadFile(chunkName)
	if err != nil {
		log.Fatal("error reading chunk: ", err)
	}

	return data, nil
}

// Delete chunk from storage server
func Delete(chunkNames []string) error {
	for _, chunkName := range chunkNames {
		err := os.Remove(chunkName)
		if err != nil {
			return err
		}
	}
	return nil
}
