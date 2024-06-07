package main

import (
	"log"
	"net"

	pb "github.com/saladin2098/month3/lesson11/public_voting/genproto"
	"github.com/saladin2098/month3/lesson11/public_voting/service"
	"github.com/saladin2098/month3/lesson11/public_voting/storage/postgres"
	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	liss, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPartyServiceServer(s,service.NewPartyService(db))
	pb.RegisterPublicServiceServer(s,service.NewPublicService(db))
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err!= nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
