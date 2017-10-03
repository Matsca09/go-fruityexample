.PHONY: all proto server client

proto:
	protoc -I . fruit/fruit.proto --go_out=plugins=grpc:.

server:
	go build -o grpc-server server/server.go 

client:
	go build -o grpc-client client/client.go 

clean:
	rm grpc-server grpc-client fruit/fruit.pb.go

all: proto server client
