package main

import (
	"flag"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/showiot/camera/gateway"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/inits/redis"
	"github.com/showiot/camera/pkg/v1/camera_messages"
	"github.com/showiot/camera/pkg/v1/cameras"
	"github.com/showiot/camera/pkg/v1/feedback"
	"github.com/showiot/camera/pkg/v1/users"
	"github.com/showiot/camera/proto/camera_messages_pb"
	"github.com/showiot/camera/proto/cameras_pb"
	"github.com/showiot/camera/proto/feedback_pb"
	"github.com/showiot/camera/proto/users_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var configPath = ""

func init()  {
	flag.StringVar(&configPath, "config_path", "configs/config.yaml", "config path")
}
func main() {
	flag.Parse()
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
	logger.GetLogger().WithError(err).Error()
	exit()
}
func exit() {
	_ = psql.Close()
	_ = redis.Close()
}
func runGrpcServer(exitChan chan error) {
	port := config.Conf.GrpcServerConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		exitChan <- err
	}
	sev := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			utils.UnaryRecover, utils.UnaryValidate,
		)),
	)
	users_pb.RegisterUserServiceServer(sev, users.NewUserServiceImpl())
	feedback_pb.RegisterFeedbackServiceServer(sev, feedback.NewFeedBackImpl())
	cameras_pb.RegisterCameraServiceServer(sev, cameras.NewCamerasImpl())
	camera_messages_pb.RegisterCameraMessageServiceServer(sev, camera_messages.NewCameraMessagesImpl())
	log.Println("Serving gRPC on port:", port)
	if err := sev.Serve(lis); err != nil {
		exitChan <- err
	}
}
