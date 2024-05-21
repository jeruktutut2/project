package setup

import (
	"log"
	"net"
	grpcuser "project-user/grpc"
	pbuser "project-user/grpc/pb/api/v1/user"
	"time"

	"google.golang.org/grpc"
)

func NewGrpcSetup(serviceSetup *ServiceSetup, port string) (grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(time.Now().String(), "error when listen: ", err)
	}
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024), // Set max receive message size (1MB in this example)
	}
	grpcServer = grpc.NewServer(opts...)
	userGrpcService := grpcuser.NewUserGrpcService(serviceSetup.UserService)
	pbuser.RegisterUserServiceServer(grpcServer, userGrpcService)

	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalln(time.Now().String(), "error when server grpc:", err)
		}
	}()
	println(time.Now().String(), "grpc: started on port", port)
	return
}

func StopGrpc(grpcServer *grpc.Server) {
	grpcServer.GracefulStop()
	println(time.Now().String(), "grpc: stopped")
}
