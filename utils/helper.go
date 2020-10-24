package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/dgrijalva/jwt-go"
	"github.com/showiot/camera/inits/config"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// 生成hash密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// 校验密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var defaultSeed = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numericSend = []rune("0123456789")
/**
 * 生成随机字符串
 * n:长度
 * randomType：字符串类型（1：字符，2数字）
 */
func GenRandomString(n, randomType int) string {
	var seed []rune
	switch randomType {
	case 1:
		seed = defaultSeed
	case 2:
		seed = numericSend
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
	if response.Code != "OK" {
		err = fmt.Errorf("send sms failed")
	}
	return
}

// 校验验证码

const (
	// access token过期时间
	ATokenExpireDuration = time.Second * 60 * 24 * 30
	// refresh token过期时间
	RTokenExpireDuration = time.Second * 3
)

// 密钥
var Secret = []byte("亚索主E，QEQ")

type Claims struct {
	UserId int32 `json:"user_id"`
	jwt.StandardClaims
}
// 生成jwt token
func GenToken(userId int32) (aToken string, err error) {
	c := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ATokenExpireDuration).Unix(),
		},
	}
	// 生成access token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(Secret)
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
