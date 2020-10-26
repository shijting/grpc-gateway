package users

import (
	"context"
	"fmt"
	go_redis "github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/inits/redis"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/users_pb"
	"github.com/showiot/camera/utils"
	"google.golang.org/grpc/codes"
	"sync"
	"time"
)

type userServiceImpl struct {
	sync.Mutex
}

func NewUserServiceImpl() *userServiceImpl {
	return &userServiceImpl{}
}

// 获取验证码
func (u *userServiceImpl) GetCode(ctx context.Context, req *users_pb.GetCodeRequest) (*empty.Empty, error) {
	//funk.Contains([]{}, req.GetCodeType())
	if req.GetCodeType() != users_pb.CodeType_CODE_LOGIN && req.CodeType != users_pb.CodeType_CODE_REGISTER {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS)
	}
	codeConfig := config.Conf.CodeConfig
	fmt.Printf("%#v\n", codeConfig)
	// retry key: recode-codeType-phoneNumber
	retryKey := fmt.Sprintf("recode-%d-%s", req.GetCodeType(), req.GetPhoneNumber())
	_, err := redis.GetDB().Get(retryKey).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err != redis.Nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS, "验证码已发送，请注意查收!")
	}
	var codeString string
	// ttl key: code-codeType-phoneNumber
	key := fmt.Sprintf("code-%d-%s", req.GetCodeType(), req.GetPhoneNumber())
	cacheCode, err := redis.GetDB().Get(key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err == redis.Nil { // 随机验证码
		codeString = utils.GenRandomString(6, 2)
	} else {
		codeString = cacheCode
	}
	// 发送验证码短信
	if err := utils.SendSms(req.GetPhoneNumber(), codeString); err != nil {
		return nil, proto.Errorg(codes.Internal, proto.Error_ERR_INTERNAL_SERVER, "验证码发送失败，清稍后重试")
	}
	// 再次重试时间
	var retryTtl int
	// 验证码过期时间
	var ttl int
	if req.GetCodeType() == users_pb.CodeType_CODE_REGISTER {
		ttl = codeConfig.RegisterTTL
		retryTtl = codeConfig.RegisterRetryTtl
	} else {
		ttl = codeConfig.LoginTTL
		retryTtl = codeConfig.LoginRetryTtl
	}
	//  记录验证码短信发送记录
	_, err =redis.GetDB().Pipelined(func(pipe go_redis.Pipeliner) (err error) {
		err =pipe.Set(key, codeString, time.Second*time.Duration(ttl)).Err()
		if err !=nil {
			return
		}
		err =pipe.Set(retryKey, 1, time.Second*time.Duration(retryTtl)).Err()
		return nil
	})
	// 记录日志
	return &empty.Empty{}, nil
}

// 注册
func (u *userServiceImpl) Register(ctx context.Context, req *users_pb.RegisterRequest) (*users_pb.UserResponse, error) {
	resq := &users_pb.UserResponse{}
	return resq, nil
}

// 登录
func (u *userServiceImpl) Login(ctx context.Context, req *users_pb.LoginRequest) (*users_pb.LoginResponse, error) {
	return nil, nil
}
