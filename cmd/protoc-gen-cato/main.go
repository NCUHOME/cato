package main

import (
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/src"
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
	pbResponse := new(pluginpb.CodeGeneratorResponse)
	generator := src.NewCatoGenerator(pbRequest)
	generator.Generate(pbResponse)
	output, err := proto.Marshal(pbResponse)
	if err != nil {
		log.Fatalf("[-] cato marshaling response error: %#v", err)
	}
	_, _ = os.Stdout.Write(output)
}
