package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/dgrijalva/jwt-go"
	"github.com/showiot/camera/gateway"
	"github.com/showiot/camera/inits/config"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"math/rand"
	"net"
	"strings"
	"time"
)

// 手机号码脱敏
func MaskedMobile(mobile string) string {
	return mobile[:3] + "***" + mobile[8:]
}

// 从token中解析出user信息
func GetUserInfoFromToken(ctx context.Context) (userInfo *UserTokenValue, err error) {
	token, err := GetCtxInfo(ctx, gateway.CtxToken)
	if err != nil || token == "" {
		return nil, fmt.Errorf("can not get token from context")
	}
	ut := New(config.Conf.TokenConfig.Prefix, 0)
	userToken, err := ut.Get(token)
	if err !=nil {
		return nil, err
	}
	return userToken, nil
}

func GetCtxInfo(ctx context.Context, key string) (val string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//logger.GetLogger().Warningfln("md not found")
		return
	}
	result := md.Get(key)
	if len(result) == 0 {
		return "", fmt.Errorf("cannot get %s from the context", key)
	}
	return result[0], nil
}

// GetClientIP 获取客户端ip
func GetClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if addSlice[0] == "[" {
		//本机地址
		return "127.0.0.1", nil
	}
	return addSlice[0], nil
}

// HashPassword 生成hash密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash 校验密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var defaultSeed = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numericSend = []rune("0123456789")

/**
# 生成随机字符串
# n:长度
# randomType：字符串类型（1：字符，2数字）
*/
func GenRandomString(n, randomType int) string {
	rand.Seed(time.Now().UnixNano())
	var seed []rune
	switch randomType {
	case 1:
		seed = defaultSeed
	case 2:
		seed = numericSend
	case 3:
		seed = []rune("123456789")
	default:
		seed = defaultSeed
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = seed[rand.Intn(len(seed)-1)]
	}
	return string(b)
}

type smsCode struct {
	Code string `json:"code"`
}

/**
 *	发送短信
 *	phoneNumber: 手机号码
 *	code： 验证码
 */
func SendSms(phoneNumber, code string) (err error) {
	var client *dysmsapi.Client
	client, err = dysmsapi.NewClientWithAccessKey(config.Conf.AliyunSmsConfig.RegionId, config.Conf.AliyunSmsConfig.AccessKeyId, config.Conf.AliyunSmsConfig.AccessSecret)
	if err != nil {
		return
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = config.Conf.AliyunSmsConfig.Scheme

	request.PhoneNumbers = phoneNumber
	request.SignName = config.Conf.AliyunSmsConfig.SignName
	request.TemplateCode = config.Conf.AliyunSmsConfig.TemplateCode
	smsCode := &smsCode{Code: code}
	byteCode, err := json.Marshal(smsCode)
	if err != nil {
		return err
	}
	request.TemplateParam = string(byteCode)
	var response *dysmsapi.SendSmsResponse
	response, err = client.SendSms(request)
	if  response.Code != "OK" {
		err = fmt.Errorf("send sms failed, error[message:%v, code:%v, requestId:%v]", response.Message, response.Code, response.RequestId)
	}
	return
}

const (
	// access token过期时间
	ATokenExpireDuration = time.Second * 60 * 24 * 30
	// refresh token过期时间
	RTokenExpireDuration = time.Second * 3
)

// 密钥
var Secret = []byte("亚索主E，QEQ")

type Claims struct {
	UserId uint32 `json:"user_id"`
	jwt.StandardClaims
}

// 生成jwt token
func GenToken(userId uint32) (token string, err error) {
	c := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ATokenExpireDuration).Unix(),
		},
	}
	// 生成access token
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(Secret)
	if err != nil {
		return
	}

	return
}

// 解析token
func ParseToken(tokenString string) (claims *Claims, err error) {
	var token *jwt.Token
	claims = new(Claims)
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
		return

	}
	return
}
