package servers

import (
	"encoding/json"
	"errors"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"time"
)

type accountInfo struct {
	SystemId     string `json:"systemId"`
	Password     string `json:"password"`
	RegisterTime int64  `json:"register_time"`
}

func Register(systemId string, password string) (err error) {
	//校验是否为空
	if len(systemId) == 0 || len(password) == 0 {
		return errors.New("系统ID或密码不能为空")
	}

	//判断是否被注册
	exist, err := redis.SISMEMBER(define.REDIS_KEY_ACCOUNT_LIST, systemId)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("该系统ID已被注册")
	}

	accountInfo := accountInfo{
		SystemId:     systemId,
		Password:     password,
		RegisterTime: time.Now().Unix(),
	}

	jsonBytes, _ := json.Marshal(accountInfo)

	//注册
	_, err = redis.Set(define.REDIS_PREFIX_ACCOUNT_INFO+systemId, string(jsonBytes))
	if err != nil {
		return err
	}

	_, err = redis.SetAdd(define.REDIS_KEY_ACCOUNT_LIST, systemId)
	if err != nil {
		return err
	}

	return nil
}

func Login(systemId string, password string) (token string, err error) {
	//校验是否为空
	if len(systemId) == 0 || len(password) == 0 {
		err = errors.New("系统ID或密码不能为空")
		return
	}

	//判断是否存在
	jsonValue, err := redis.Get(define.REDIS_PREFIX_ACCOUNT_INFO + systemId)
	if err != nil || len(jsonValue) == 0 {
		err = errors.New("系统ID或密码错误")
		return
	}
	accountInfo := accountInfo{}

	_ = json.Unmarshal([]byte(jsonValue), &accountInfo)
	if systemId != accountInfo.SystemId || password != accountInfo.Password {
		err = errors.New("系统ID或密码错误")
		return
	}

	//生成token
	token = util.GenUUID()

	//存到redis
	_, err = redis.SetWithSurvivalTime(define.REDIS_PREFIX_TOKEN+token, systemId, define.REDIS_KEY_SURVIVAL_SECONDS)
	if err != nil {
		return "", err
	}

	return
}
