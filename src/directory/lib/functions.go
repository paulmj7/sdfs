package lib

import (
	"log"
	"net/url"
)

type StorageServers struct {
	servers []url.URL
}

func (ss *StorageServers) Add(server string) {
	url, err := url.Parse(server)
	if err != nil {
		log.Fatal("storage server url not valid: ", err)
	}
	ss.servers = append(ss.servers, *url)
}

func (ss *StorageServers) PrintDirectory() {
	for i, s := range ss.servers {
		log.Println(i, s.String())
	}
}

type DFile struct {
	Chunks   []Chunk
	FileSize int
	Len      int
}

type Chunk struct {
	Location        url.URL
	BackupLocations []url.URL
	ChunkSize       int
}

var Lookup map[string]DFile
var Directory StorageServers
