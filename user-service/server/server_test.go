package main

import (
	"context"
	users "github.com/lhj8390/grpc-example/user-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10)
	s := grpc.NewServer()
	registerServices(s)

	go func() {
		err := startServer(s, l)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return s, l
}

func TestUserService(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop()

	bufconnDialer := func(
		ctx context.Context,
		addr string,
	) (net.Conn, error) {
		return l.Dial()
	}

	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}
	usersClient := users.NewUsersClient(client)
	resp, err := usersClient.GetUser(
		context.Background(),
		&users.UserGetRequest{
			Email: "lhj8390@naver.com",
			Id:    "foo-bar",
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	if resp.User.FirstName != "lhj8390" {
		t.Errorf(
			"Expected FirstName to be: lhj8390, Got: %s",
			resp.User.FirstName,
		)
	}
}
