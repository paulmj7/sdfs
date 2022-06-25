package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"math"
	"net/url"
	"sdfs/services/pb"
	"strconv"
)

// StorageServers struct
type StorageServers struct {
	servers []url.URL
}

// Add storage server to directory
func (ss *StorageServers) Add(server string) {
	url, err := url.Parse(server)
	if err != nil {
		log.Fatal("storage server url not valid: ", err)
	}
	ss.servers = append(ss.servers, *url)
}

// Size of StorageServers struct
func (ss *StorageServers) Size() int {
	return len(ss.servers)
}

// PrintDirectory prints contents of StorageServers
func (ss *StorageServers) PrintDirectory() {
	for i, s := range ss.servers {
		log.Println(i, s.String())
	}
}

// Create chunk mappings for file
func Create(fileName string, fileSize uint64) ([]string, error) {
	numFiles := int(math.Ceil(float64(fileSize) / 64000000))
	remainder := fileSize % 64000000
	log.Println(numFiles, remainder)
	chunks := []Chunk{}
	h := sha1.New()
	for i := 0; i < numFiles; i++ {
		var chunkSize uint64 = 64000000
		if i+1 == numFiles {
			chunkSize = remainder
		}
		location := location()
		io.WriteString(h, fileName+strconv.Itoa(i))
		name := hex.EncodeToString(h.Sum(nil))[:10]
		chunk := Chunk{Name: name, Location: location, BackupLocations: nil, ChunkSize: chunkSize}
		chunks = append(chunks, chunk)
	}
	f := DFile{Chunks: chunks, FileSize: fileSize, Len: numFiles}
	Lookup[fileName] = f
	locations := []string{}
	for _, chunk := range f.Chunks {
		location := chunk.Location.Hostname() + ":" + chunk.Location.Port()
		locations = append(locations, location)
	}
	return locations, nil
}

func location() url.URL {
	if Directory.Size() == 1 {
		location := Directory.servers[0]
		return location
	}
	panic("error getting location")
}

// Search file in directory
func Search(fileName string) ([]*pb.ReadChunk, error) {
	f, exists := Lookup[fileName]
	if !exists {
		log.Fatal("error finding file")
		return nil, errors.New("error finding dfile")
	}
	chunks := []*pb.ReadChunk{}
	for _, chunk := range f.Chunks {
		name := chunk.Name
		location := chunk.Location.Hostname() + ":" + chunk.Location.Port()
		chunks = append(chunks, &pb.ReadChunk{Name: name, Location: location})
	}

	return chunks, nil
}

// Ls lists all files in directory
func Ls() ([]string, error) {
	names := make([]string, len(Lookup))
	i := 0
	for key := range Lookup {
		names[i] = key
		i++
	}

	return names, nil
}

// DFile struct
type DFile struct {
	Chunks   []Chunk
	FileSize uint64
	Len      int
}

// Chunk struct
type Chunk struct {
	Name            string
	Location        url.URL
	BackupLocations []url.URL
	ChunkSize       uint64
}

// Lookup map
var Lookup = map[string]DFile{}

// Directory object
var Directory StorageServers
