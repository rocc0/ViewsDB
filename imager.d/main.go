package main

import (
	"log"
	"net"

	pb "./pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func init() {
	//INIT config vars, MUST be first in init statement!!!
	if err := config.getConf(); err != nil {
		log.Fatalf("Error when parsing config: %v", err)
	}
	//INIT config vars, MUST be first in init statement!!!

	cli, err := NewConsulClient(config.Consul)
	if err != nil {
		log.Fatal(err)
	}
	cli.Register("imager", 345)
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
