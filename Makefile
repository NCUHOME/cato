

build-proto:
	protoc -I=../ --go_out=../../../ ../cato/proto/*.proto


install:
	go install ./cmd/protoc-gen-cato