build-proto:
	protoc -I=../ --go_out=../../../ ../cato/proto/*.proto

link-local:
	ln -s "$(pwd)" "$(dirname $(dirname $(which protoc)))"

install:
	go install ./cmd/protoc-gen-cato