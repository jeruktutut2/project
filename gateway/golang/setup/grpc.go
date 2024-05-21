package setup

import (
	"log"
	"time"

	pbuser "gateway/protofiles/pb/api/v1/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClientConnection(host string) (*grpc.ClientConn, pbuser.UserServiceClient) {
	println(time.Now().String(), "user client connection: connecting to", host)
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// defer conn.Close()
	if err != nil {
		log.Fatalln("user client connection: error when new client", err)
	}
	println(time.Now().String(), "user client connection: connected to", host)
	return conn, pbuser.NewUserServiceClient(conn)
}

func CloseUserClientConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Fatalln(time.Now().String(), "error when close user client connection:", err)
	}
	println(time.Now().String(), "user client connection: closed properly")
}
