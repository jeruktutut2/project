package setup

import (
	"log"
	"net"
	grpcuser "project-user/grpc"
	pbuser "project-user/grpc/pb/api/v1/user"
	"project-user/interceptor"
	"time"

	"google.golang.org/grpc"
)

func NewUserGrpcSetup(serviceSetup *ServiceSetup, host string) (grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalln(time.Now().String(), "error when listen: ", err)
	}
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024), // Set max receive message size (1MB in this example)
		grpc.UnaryInterceptor(interceptor.SetLog),
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
	println(time.Now().String(), "grpc user: started on host", host)
	return
}

func StopUserGrpc(grpcServer *grpc.Server) {
	grpcServer.GracefulStop()
	println(time.Now().String(), "grpc user: stopped")
}
