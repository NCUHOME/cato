package main

import (
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protoInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("[-] cato read data from stdin: %#v", err)
	}
	pbRequest := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(protoInput, pbRequest); err != nil {
		log.Fatalf("[-] cato unmarshal pbRequest data: %#v", err)
	}

}
