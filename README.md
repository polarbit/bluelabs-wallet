# go-bootcamp


# Install 
go get github.com/AsynkronIT/protoactor-go
go get github.com/spf13/cobra

# temp
https://github.com/protocolbuffers/protobuf
https://github.com/protocolbuffers/protobuf-go
https://grpc.io/docs/languages/go/quickstart/
https://github.com/gogo/protobuf


$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

go install github.com/gogo/protobuf/protoc-gen-gofast
go install  github.com/gogo/protobuf/protoc-gen-gogoslick
go install github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograinv2@dev

#
protoc --gogoslick_out=. ./worker/messages/protos.proto
protoc --gograinv2_out=. ./worker/messages/protos.proto


# Bunu kullanmaya gerek yok gibi
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative  protos.proto

