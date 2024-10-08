package main

import (
	"os"

	"github.com/guoming0000/protoc-gen-go-gin/cmd/protoc-gen-openapi/generator"
)

func main() {
	if len(os.Args) < 3 {
		return
	}
	fullPath := os.Args[1]
	newFilePath := os.Args[2]

	generator.RecreateProto(fullPath, newFilePath)
}
