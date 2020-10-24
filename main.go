package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/showiot/camera/gateway"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/pkg/v1/users"
	"github.com/showiot/camera/proto/users_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)
const configPath = "configs/config.yaml"
func main() {
	if err :=config.Init(configPath);err !=nil {
		log.Fatal(err)
	}
	exitChan := make(chan error)
	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		exitChan <- fmt.Errorf("%s", <-signalChan)
	}()
	go runGrpcServer(exitChan)
	go gateway.Run(exitChan)
	err := <-exitChan
	log.Println(err)
}
func runGrpcServer(exitChan chan error) {
	port := config.Conf.GrpcServerConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		exitChan <- err
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			utils.UnaryRecover, utils.UnaryValidate,
		)),
	)
	users_pb.RegisterUserServiceServer(s, users.NewUserServiceImpl())
	log.Println("Serving gRPC on port:", port)
	if err := s.Serve(lis); err != nil {
		exitChan <- err
	}
}
