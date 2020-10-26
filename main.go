package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/showiot/camera/gateway"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/inits/redis"
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
	if err := config.Init(configPath); err != nil {
		log.Fatal(err)
	}
	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}
	if err := psql.Init(); err != nil {
		log.Fatal(err)
	}
	if err := redis.Init(); err != nil {
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
	exit()
}
func exit() {
	_ = psql.Close()
	_ = redis.Close()
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
