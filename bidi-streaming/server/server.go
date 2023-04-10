package main

import (
	"fmt"
	svc "github.com/lhj8390/grpc-example/bidi-streaming/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

type userService struct {
	svc.UnimplementedUsersServer
}

func (s *userService) GetHelp(
	stream svc.Users_GetHelpServer,
) error {
	log.Println("Client connected")

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Request received: %s\n", request.Request)
		response := svc.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}

	log.Printf("Client disconnected")
	return nil
}

func registerServices(s *grpc.Server) {
	svc.RegisterUsersServer(s, &userService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	registerServices(s)
	log.Fatal(startServer(s, lis))
}
