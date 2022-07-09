package logger

import (
	"log"
	"os"
)

func DebugLog(msg string) {
	if os.Getenv("DEBUG") == "true" {
		log.Println(msg)
	}
}
