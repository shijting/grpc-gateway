package users

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10"
	go_redis "github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/inits/logger"
	"github.com/showiot/camera/inits/psql"
	"github.com/showiot/camera/inits/redis"
	"github.com/showiot/camera/models"
	"github.com/showiot/camera/proto"
	"github.com/showiot/camera/proto/users_pb"
	"github.com/showiot/camera/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"regexp"
	"time"
)

var (
	ErrorRegisterFailed       = errors.New("注册失败")
	ErrorUserNotExist         = errors.New("用户不存在")
	ErrorInvalidCaptcha       = errors.New("验证码错误")
	ErrorLoginFailed          = errors.New("登录失败")
	ErrorUpdateUserInfoFailed = errors.New("更新用户资料失败")
)

// 1 megabyte
const maxImageSize = 1 << 20

type impl struct {
	captchaKey      string
	captchaRetryKey string
}

func NewUserServiceImpl() *impl {
	return &impl{
		captchaKey:      "code-%s",
		captchaRetryKey: "recode-%s",
	}
}

func (u *impl) getCaptchaRetryKey(phoneNumber string) string {
	return fmt.Sprintf(u.captchaRetryKey, phoneNumber)
}

// getCaptchaKey
func (u *impl) getCaptchaKey(phoneNumber string) string {
	return fmt.Sprintf(u.captchaKey, phoneNumber)
}

// GetCaptcha 获取验证码
func (u *impl) GetCaptcha(ctx context.Context, req *users_pb.GetCaptchaRequest) (*empty.Empty, error) {
	//funk.Contains([]{}, req.GetCodeType())
	codeConfig := config.Conf.CodeConfig
	retryKey := u.getCaptchaRetryKey(req.GetPhoneNumber())
	_, err := redis.GetDB().Get(retryKey).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err != redis.Nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS, "验证码已发送，请注意查收!")
	}
	var codeString string
	key := u.getCaptchaKey(req.GetPhoneNumber())
	cacheCode, err := redis.GetDB().Get(key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err == redis.Nil { // 随机验证码
		codeString = utils.GenRandomString(6, 2)
	} else {
		codeString = cacheCode
	}
	fields := logrus.Fields{
		"phone_number": req.GetPhoneNumber(),
	}
	// 发送验证码短信
	if err := utils.SendSms(req.GetPhoneNumber(), codeString); err != nil {
		logger.GetLogger().WithFields(fields).WithError(err).Error("send sms failed")
		return nil, proto.Errorg(codes.Internal, proto.Error_ERR_INTERNAL_SERVER, "验证码发送失败，清稍后重试")
	}
	// 再次重试时间
	var retryTtl int
	// 验证码过期时间
	var ttl int
	ttl = codeConfig.LoginTTL
	retryTtl = codeConfig.LoginRetryTtl
	//  记录验证码短信发送记录
	_, err = redis.GetDB().Pipelined(func(pipe go_redis.Pipeliner) (err error) {
		err = pipe.Set(key, codeString, time.Second*time.Duration(ttl)).Err()
		if err != nil {
			logger.GetLogger().WithFields(fields).WithError(err).Error("")
		}
		err = pipe.Set(retryKey, 1, time.Second*time.Duration(retryTtl)).Err()
		if err != nil {
			logger.GetLogger().WithFields(fields).WithError(err).Error("")
		}
		return nil
	})
	// 记录短信发送日志
	logger.GetLogger().WithFields(fields).Info("send sms successful")
	return &empty.Empty{}, nil
}

// VerifyCaptcha 校验验证码
func (u *impl) VerifyCaptcha(ctx context.Context, req *users_pb.VerifyCaptchaRequest) (*users_pb.VerifyCaptchaResponse, error) {
	var (
		key, captcha string
		err          error
		verified     = true
	)
	fields := logrus.Fields{
		"func":         "VerifyCaptcha",
		"phone_number": req.GetPhoneNumber(),
	}
	key = u.getCaptchaKey(req.GetPhoneNumber())
	// 获取验证码失败
	if captcha, err = redis.GetDB().Get(key).Result(); err != nil {
		logger.GetLogger().WithFields(fields).WithError(err).Error("")
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INVALID_CAPTCHA, ErrorInvalidCaptcha.Error())
	}
	// 验证码不正确
	if req.GetCaptcha() != captcha {
		verified = false
	}

	return &users_pb.VerifyCaptchaResponse{Verified: verified}, nil
}

// Login 登录
func (u *impl) Login(ctx context.Context, req *users_pb.LoginRequest) (*users_pb.LoginResponse, error) {
	var (
		dbUser  *models.User
		user    *models.User
		result  pg.Result
		err     error
		captcha string
		token   uuid.UUID
		ip      string
	)

	fields := logrus.Fields{
		"phone_number": req.GetPhoneNumber(),
	}
	key := u.getCaptchaKey(req.GetPhoneNumber())
	// 获取验证码失败
	if captcha, err = redis.GetDB().Get(key).Result(); err != nil {
		logger.GetLogger().WithFields(fields).WithError(err).Error("get cache from redis failed key:", key)
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER, ErrorInvalidCaptcha.Error())
	}
	// 验证码不正确
	if req.GetCaptcha() != captcha {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_CAPTCHA, ErrorInvalidCaptcha.Error())
	}
	dbUser = new(models.User)
	err = psql.GetDB().Model(dbUser).Where("phone_number = ?", req.GetPhoneNumber()).First()
	if err != nil && err != pg.ErrNoRows {
		logger.GetLogger().WithFields(fields).WithError(err).Error("login error")
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_OPERATION_FAILED, ErrorLoginFailed.Error())
	}
	// 用户不存在，先注册
	if dbUser.Id == 0 {
		user = &models.User{
			PhoneNumber: req.GetPhoneNumber(),
		}
		result, err = psql.GetDB().Model(user).Insert()
		if err != nil || result.RowsAffected() == 0 {
			logger.GetLogger().WithFields(fields).Error(err.Error())
			return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_OPERATION_FAILED, ErrorRegisterFailed.Error())
		}
	}
	go func() {
		// 更新用户信息
		if ip, err = utils.GetClientIP(ctx); err != nil {
			logger.GetLogger().WithFields(fields).Error("get remote addr err: ", err.Error())
		}
		user = &models.User{
			Id:          dbUser.Id,
			LastLoginAt: time.Now(),
			LastLoginIp: ip,
			UpdatedAt:   time.Now(),
		}
		_, err = psql.GetDB().Model(user).WherePK().Column("last_login_at", "last_login_ip", "updated_at").Update()
		if err != nil {
			logger.GetLogger().WithFields(fields).WithError(err).Error("")
		}
	}()

	// 获取token
	userToken := utils.New(config.Conf.TokenConfig.Prefix, config.Conf.TokenConfig.Expire)
	userTokenVal := utils.NewValue(dbUser.Id)
	if token, err = uuid.NewRandom(); err != nil {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_GENERATE_TOKEN_FAILED)
	}
	if err = userToken.Set(fmt.Sprintf("%s", token), userTokenVal); err != nil {
		logger.GetLogger().WithFields(fields).WithError(err).Error("login error")
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_ARGS, ErrorLoginFailed.Error())
	}
	// 登录成功
	return &users_pb.LoginResponse{Token: fmt.Sprintf("%s", token)}, nil
}

// 更新用户信息
func (u *impl) UpdateUserInfo(ctx context.Context, req *users_pb.UpdateUserInfoRequest) (*users_pb.UpdateUserInfoResponse, error) {
	var (
		tokenVal *utils.UserTokenValue
		err      error
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	data := &models.User{Id: tokenVal.UserID, Nickname: req.GetNickname()}
	_, err = psql.GetDB().Model(data).WherePK().Column("nickname").Update()
	if err != nil {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_USER_NOT_EXIST, ErrorUpdateUserInfoFailed.Error())
	}
	return &users_pb.UpdateUserInfoResponse{Nickname: req.GetNickname()}, nil
}

func (u *impl) ShowUser(ctx context.Context, _ *empty.Empty) (resp *users_pb.User, err error) {
	var (
		result   models.User
		tokenVal *utils.UserTokenValue
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}
	err = psql.GetDB().Model(&result).Where("id = ?", tokenVal.UserID).First()
	if err != nil {
		return nil, proto.Errorg(codes.NotFound, proto.Error_ERR_USER_NOT_EXIST, ErrorUserNotExist.Error())
	}
	resp = new(users_pb.User)
	copier.Copy(resp, &result)
	return resp, nil
}


func (u *impl) Upload(ctx context.Context, req *users_pb.UploadAvatarRequest) (*users_pb.UploadAvatarResponse, error) {
	var (
		user   *models.User
		tokenVal *utils.UserTokenValue
		imageUrl string
		err      error
		db = psql.GetDB()
	)
	if tokenVal, err = utils.GetUserInfoFromToken(ctx); err != nil {
		return nil, proto.Errorg(codes.InvalidArgument, proto.Error_ERR_INVALID_TOKEN)
	}

	if imageUrl, err = WriteFile("uploads", req.GetAvatar()); err != nil {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	user = &models.User{
		Id:        tokenVal.UserID,
		Avatar:    imageUrl,
		UpdatedAt: time.Now(),
	}
	if _, err = db.Model(user).WherePK().Column("avatar", "updated_at").Update(); err != nil {
		return nil, proto.Errorg(codes.Unknown, proto.Error_ERR_INTERNAL_SERVER)
	}
	return &users_pb.UploadAvatarResponse{ImageUrl: imageUrl}, nil
}

func WriteFile(path string, base64ImageContent string) (string, error) {
	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64ImageContent)
	if !b {
		return "", fmt.Errorf("only supported image file")
	}

	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64ImageContent), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取
	fmt.Println(fileType)
	base64Str := re.ReplaceAllString(base64ImageContent, "")

	UUID, err :=uuid.NewRandom()
	if err != nil {
		return "", err
	}
	var file string = path + "/" + fmt.Sprintf("%s", UUID) + "." + fileType
	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err = ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		return "", err
	}

	return file, nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return proto.Errorg(codes.Canceled, proto.Error_ERR_INTERNAL_SERVER, "request is canceled")
	case context.DeadlineExceeded:
		return proto.Errorg(codes.DeadlineExceeded, proto.Error_ERR_INTERNAL_SERVER, "deadline is exceeded")
	default:
		return nil
	}
}
