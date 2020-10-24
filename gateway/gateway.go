package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/middlewares"
	"github.com/showiot/camera/proto/users_pb"
	_ "github.com/showiot/camera/statik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"mime"
	"net/http"
	"strings"
)

func Run(exitChan chan error) {
	errorOption := runtime.WithErrorHandler(gatewayErrorHandler)
	mux := runtime.NewServeMux(errorOption)
	ctx := context.Background()
	serverPort := config.Conf.GrpcServerConfig.Port
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf(":%d", serverPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		exitChan <- err
	}
	if err := users_pb.RegisterUserServiceHandler(ctx, mux, conn); err != nil {
		exitChan <- err
	}
	openAPIHandler := getOpenAPIHandler(exitChan)

	m := NewMiddleware()
	// 使用日志中间件
	m.Use(middlewares.LoggerMiddleware)
	// 权限验证
	m.Use(middlewares.JWTAuthMiddleware)
	httpPort := config.Conf.GrpcGwServerConfig.Port
	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", httpPort),
		Handler: m.Add(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				mux.ServeHTTP(w, r)
				return
			}
			openAPIHandler.ServeHTTP(w, r)
		})),
	}
	log.Println("Serving gRPC-gateway on port:", httpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		exitChan <- err
	}
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
	Code uint32 `json:"code"`
	Message string `json:"message"`
}
// 自定义错误处理
func gatewayErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "application/json")
	respData, _ :=json.Marshal(errorResponse{
		Code:    uint32(st.Code()),
		Message: st.Message(),
	})
	w.Write(respData)
}
