package servers

import (
	"encoding/json"
	"errors"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"time"
)

type accountInfo struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RegisterTime int64  `json:"register_time"`
}

func Register(username string, password string) (err error) {
	//校验是否为空
	if len(username) == 0 || len(password) == 0 {
		return errors.New("用户名或密码不能为空")
	}

	//判断是否被注册
	exist, err := redis.SISMEMBER(define.REDIS_KEY_ACCOUNT_LIST, username)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("该用户账号已被注册")
	}

	accountInfo := accountInfo{
		Username:     username,
		Password:     password,
		RegisterTime: time.Now().Unix(),
	}

	jsonBytes, _ := json.Marshal(accountInfo)

	//注册
	_, err = redis.Set(define.REDIS_PREFIX_ACCOUNT_INFO+username, string(jsonBytes))
	if err != nil {
		return err
	}

	_, err = redis.SetAdd(define.REDIS_KEY_ACCOUNT_LIST, username)
	if err != nil {
		return err
	}

	return nil
}
