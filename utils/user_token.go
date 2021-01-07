package utils

import (
	"encoding/json"
	"fmt"
	go_redis "github.com/go-redis/redis"
	"github.com/showiot/camera/inits/redis"
	log "github.com/sirupsen/logrus"
	"time"
)


type UserTokenValue struct {
	UserID    uint32
	Type      int
	ExpiredAt time.Time
}

func NewValue(userId uint32) *UserTokenValue {
	return &UserTokenValue{
		UserID: userId,
	}
}

// 用户令牌
type UserToken struct {
	expires      int
	tokenListKey string
	tokenKey     string
}

func New(prefix string, expires int) *UserToken {
	return &UserToken{
		expires:      expires,
		tokenListKey: prefix + ":user:token:list:%d",
		tokenKey:     prefix + ":user:token",
	}
}

func (ut *UserToken) listKey(userId uint32) string {
	return fmt.Sprintf(ut.tokenListKey, userId)
}

func (ut *UserToken) key(token string) string {
	return ut.tokenKey
}

func (ut *UserToken) Set(token string, val *UserTokenValue) error {
	rdb := redis.GetDB()
	listKey := ut.listKey(val.UserID)
	key := ut.key(token)

	if ut.expires > 0 {
		val.ExpiredAt = time.Now().Add(time.Second * time.Duration(ut.expires))
	}
	body, err := json.Marshal(val)

	if err != nil {
		return err
	}
	_, err = rdb.LPush(listKey, token).Result()

	if err != nil {
		return err
	}
	_, err = rdb.HSet(key, token, body).Result()
	//_, err = rdb.Do("HSET", key, token, body).Result()
	return err
}

// 删除用户的历史token
// 比如当修改密码的时候，所有设备重新登录
func (ut *UserToken) Delete(userId uint32) error {
	rdb := redis.GetDB()

	listKey := ut.listKey(userId)
	tokens, err := rdb.LRange(listKey, 0, -1).Result()
	if err != nil {
		return err
	}
	// TODO: 批量操作
	rdb.Del(listKey)
	for _, token := range tokens {
		key := ut.key(token)
		rdb.HDel(key, token)
	}
	return err
}

func (ut *UserToken) Get(token string) (val *UserTokenValue, err error) {
	rdb := redis.GetDB()

	val = &UserTokenValue{}
	key := ut.key(token)
	body, err := rdb.HGet(key, token).Result()
	if err != nil {
		return val, err
	}
	if err == go_redis.Nil {
		return val, fmt.Errorf("invalid token")
	}
	err = json.Unmarshal([]byte(body), val)
	return
}

// 获取token并且检测token是否有效
func (ut *UserToken) Check(token string) (val *UserTokenValue, ok bool) {
	value, err := ut.Get(token)
	if err != nil {
		log.WithError(err).WithField("token", token).Warnf("get token error")
		return nil, false
	}
	if time.Now().Sub(value.ExpiredAt) > 0 {
		return nil, false
	}
	return value, true
}

// 删除单个token
func (ut *UserToken) DeleteToken(userId uint32, token string) error {
	rdb := redis.GetDB()
	listKey := ut.listKey(userId)
	key := ut.key(token)
	_, err := rdb.LRem(listKey, 0, token).Result()
	if err != nil {
		return err
	}
	_, err = rdb.HDel(key, token).Result()
	return err
}
