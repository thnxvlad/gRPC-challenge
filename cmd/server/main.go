package main

import (
	"challenge/pkg/api"
	"challenge/pkg/server"
	"challenge/util"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &server.GRPCServer{}
	api.RegisterChallengeServiceServer(s, srv)

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	l, err := net.Listen("tcp", config.Port)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
