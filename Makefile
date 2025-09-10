



build-proto:
	protoc -I=./proto --go_out=../ ./proto/*.proto