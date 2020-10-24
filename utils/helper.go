package utils
import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// 生成密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// 校验密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
// 发送短信

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
