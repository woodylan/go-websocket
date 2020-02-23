package servers

import (
	"encoding/json"
	"errors"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"time"
)

type accountInfo struct {
	SystemId     string `json:"systemId"`
	RegisterTime int64  `json:"register_time"`
}

func Register(systemId string) (err error) {
	//校验是否为空
	if len(systemId) == 0 {
		return errors.New("系统ID不能为空")
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
