package main

import (
	"log"
	"net"

	pb "./imagegrpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func init() {
	cli, err := NewConsulClient("http://192.168.99.100:8500")
	if err != nil {
		log.Fatal(err)
	}
	cli.Register("kek", 345)
	if err := config.getConf(); err != nil {
		log.Fatalf("Error when parsing config: %v", err)
	}

	if err := mgoConnect(); err != nil {
		log.Fatalf("Error initializing mongo: %v\n", err)
	}
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterImagerServer(s, &server{})
	log.Print("Image server started!")
	s.Serve(lis)

}
