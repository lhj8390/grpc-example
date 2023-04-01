package main

import (
	"context"
	"errors"
	users "github.com/lhj8390/grpc-example/user-service/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

type UserService struct {
	users.UnimplementedUsersServer
}

func (s *UserService) GetUser(
	ctx context.Context,
	in *users.UserGetRequest,
) (*users.UserGetReply, error) {
	log.Printf(
		"Received request for user with Email: %s Id: %s\n",
		in.Email,
		in.Id,
	)

	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}

	// 더미 User 객체 생성
	u := users.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func registerServices(s *grpc.Server) {
	users.RegisterUsersServer(s, &UserService{})
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
