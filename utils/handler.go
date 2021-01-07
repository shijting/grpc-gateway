package utils

import (
	"context"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// 拦截器 接管panic
func UnaryRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().WithField("type", "panic").Error(r)
			//err = status.Error(codes.Internal, "系统错误，请稍后重试")
			err = proto.Errorg(codes.Internal, proto.Error_ERR_INTERNAL_SERVER)
		}
	}()
	return handler(ctx, req)
}

type validator interface {
	Validate() error
}
// 拦截器 统一参数验证
func UnaryValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	if r,ok := req.(validator);ok {
		if validateErr := r.Validate();validateErr !=nil {
			err = proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS, validateErr.Error())
			return
		}
	}
	return handler(ctx, req)
}
