package users

import (
	"context"
	"fmt"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/users_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc/codes"
	"sync"
)

type userServiceImpl struct {
	sync.Mutex
}

func NewUserServiceImpl() *userServiceImpl {
	return &userServiceImpl{}
}
// 获取验证码
func (u *userServiceImpl) GetCode (ctx context.Context, req *users_pb.GetCodeRequest) (*users_pb.GetCodeResponse, error){
	codeString := utils.GenRandomString(6, 2)

	if err :=utils.SendSms(req.GetPhoneNumber(), codeString);err !=nil {
		return nil, proto.Errorg(codes.Internal, proto.Error_ERR_INTERNAL_SERVER, "验证码发送失败，清稍后重试")
	}
	// todo 记录发送短信记录
	// todo redis记录验证码过期时间
	return &users_pb.GetCodeResponse{}, nil
}
// 注册
func (u *userServiceImpl) Register(ctx context.Context, req *users_pb.UserRegisterRequest) (*users_pb.UserResponse, error) {
	err :=utils.SendSms(req.GetPhoneNumber(), "888888")
	if err !=nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS, "发送验证码失败")
	}
	fmt.Println(utils.GenRandomString(6, 2))
	resq :=&users_pb.UserResponse{
		Code:    200,
		Message: req.GetPhoneNumber(),
		Details: nil,
	}
	return resq, nil
}
// 登录
func (u *userServiceImpl)Login(ctx context.Context, req *users_pb.UserLoginRequest) (*users_pb.UserResponse, error)  {
	return nil, nil
}