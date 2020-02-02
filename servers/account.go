package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
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

func Login(username string, password string) (token string, err error) {
	//校验是否为空
	if len(username) == 0 || len(password) == 0 {
		err = errors.New("用户名或密码不能为空")
		return
	}

	//判断是否存在
	jsonValue, err := redis.Get(define.REDIS_PREFIX_ACCOUNT_INFO + username)
	if err != nil || len(jsonValue) == 0 {
		err = errors.New("用户名或密码错误")
		return
	}
	accountInfo := accountInfo{}

	_ = json.Unmarshal([]byte(jsonValue), &accountInfo)
	fmt.Println(accountInfo)
	if username != accountInfo.Username || password != accountInfo.Password {
		err = errors.New("用户名或密码错误")
		return
	}

	//生成token
	token = util.GenUUID()

	//存到redis
	_, err = redis.SetWithSurvivalTime(define.REDIS_PREFIX_TOKEN+token, username, define.REDIS_KEY_SURVIVAL_SECONDS)
	if err != nil {
		return "", err
	}

	return
}
