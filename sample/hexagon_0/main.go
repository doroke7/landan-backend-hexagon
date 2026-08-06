package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	container "example/container"
	pb "example/pb/grpc"
)

func main() {

	oContainer, err := container.InitContainer(
		"root:root@tcp(localhost:3306)/example?parseTime=true",
		"amqp://guest:guest@localhost:5672/",
		"localhost:9090",
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := oContainer.ConsumerUserHandler.Start(context.Background()); err != nil {
			log.Printf("consumer stopped: %v", err)
		}
	}()

	go func() {
		if err := oContainer.ClientUserHandler.Start(context.Background()); err != nil {
			log.Printf("user stream client stopped: %v", err)
		}
	}()

	oGrpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(oGrpcServer, oContainer.GrpcUserHandler)

	oListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(oGrpcServer.Serve(oListener))
	}()

	http.HandleFunc("/user/create", oContainer.HttpUserHandler.CreateUser)
	http.HandleFunc("/user/get", oContainer.HttpUserHandler.GetUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
