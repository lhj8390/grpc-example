package main

import (
	"context"
	"fmt"
	users "github.com/lhj8390/grpc-example/user-service/service"
	"google.golang.org/grpc"
	"log"
	"os"
)

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
}

func getUserServiceClient(conn *grpc.ClientConn) users.UsersClient {
	return users.NewUsersClient(conn)
}

func getUser(
	client users.UsersClient,
	u *users.UserGetRequest,
) (*users.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(
			"Must specify a gRPC server address",
		)
	}

	conn, err := setupGrpcConnection(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	result, err := getUser(
		c,
		&users.UserGetRequest{Email: "lhj8390@naver.com"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(
		os.Stdout, "User: %s %s\n",
		result.User.FirstName,
		result.User.LastName,
	)
}
