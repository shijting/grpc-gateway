package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/middlewares"
	"github.com/showiot/camera/pkg/websocket"
	"github.com/showiot/camera/proto/camera_messages_pb"
	"github.com/showiot/camera/proto/cameras_pb"
	"github.com/showiot/camera/proto/feedback_pb"
	"github.com/showiot/camera/proto/users_pb"
	_ "github.com/showiot/camera/statik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"mime"
	"net/http"
	"strings"
)

const (
	CtxToken = "X-Camera-Token"
)

func Run(exitChan chan error) {
	errorOption := runtime.WithErrorHandler(gatewayErrorHandler)
	metadataOptions := runtime.WithMetadata(func(c context.Context, req *http.Request) metadata.MD {
		return metadata.Pairs(
			CtxToken, req.Header.Get(CtxToken),
		)
	})
	marshalOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true}})
	mux := runtime.NewServeMux(errorOption, metadataOptions, marshalOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverPort := config.Conf.GrpcServerConfig.Port
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", serverPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		exitChan <- err
	}
	registerConfig := []struct {
		name         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
	}{
		{
			name:         "user",
			registerFunc: users_pb.RegisterUserServiceHandler,
		},
		{
			name:         "feedback",
			registerFunc: feedback_pb.RegisterFeedbackServiceHandler,
		},
		{
			name:         "camera",
			registerFunc: cameras_pb.RegisterCameraServiceHandler,
		},
		{
			name:         "cameraMessage",
			registerFunc: camera_messages_pb.RegisterCameraMessageServiceHandler,
		},
	}
	for _, s := range registerConfig {
		if err := s.registerFunc(ctx, mux, conn); err != nil {
			exitChan <- err
		}
	}
	openAPIHandler := getOpenAPIHandler(exitChan)
	m := NewMiddleware()
	m.Use(middlewares.CorsMiddleware)
	// 使用日志中间件
	m.Use(middlewares.LoggerMiddleware)
	// 权限验证
	//m.Use(middlewares.JWTAuthMiddleware)
	httpPort := config.Conf.GrpcGwServerConfig.Port
	// websocket handler
	wsHandler := getWsHandler()
	httpServer := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", httpPort),
		Handler: m.Add(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// api
			if strings.HasPrefix(r.URL.Path, "/v1") {
				mux.ServeHTTP(w, r)
				return
			}
			// websocket
			if strings.HasPrefix(r.URL.Path, "/ws") {
				wsHandler.ServeHTTP(w, r)
				return
			}
			// swagger open api documents
			openAPIHandler.ServeHTTP(w, r)
		})),
	}
	log.Println("Serving gRPC-gateway on port:", httpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		exitChan <- err
	}
}

// get websocket handler
func getWsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})
}

func getOpenAPIHandler(quit chan error) http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		quit <- err
	}
	statikFS, err := fs.New()
	if err != nil {
		quit <- err
	}
	return http.FileServer(statikFS)
}

type errorResponse struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

// 自定义错误处理
func gatewayErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	respData, _ := json.Marshal(errorResponse{
		Code:    uint32(st.Code()),
		Message: st.Message(),
	})
	w.Write(respData)
}
