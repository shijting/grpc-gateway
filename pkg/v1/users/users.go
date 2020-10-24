package users

import (
	"context"
	"github.com/showiot/camera/proto/users_pb"
	"sync"
	"time"
)

type userServiceImpl struct {
	sync.Mutex
}

func NewUserServiceImpl() *userServiceImpl {
	return &userServiceImpl{}
}
func (u *userServiceImpl) Register(ctx context.Context, req *users_pb.UserRegisterRequest) (*users_pb.UserResponse, error) {
	resq :=&users_pb.UserResponse{
		Code:    200,
		Message: req.GetPhoneNumber(),
		Details: nil,
	}
	time.Sleep(1*time.Second)
	return resq, nil
}
func (u *userServiceImpl)Login(ctx context.Context, req *users_pb.UserLoginRequest) (*users_pb.UserResponse, error)  {
	return nil, nil
}