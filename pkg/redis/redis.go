package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-websocket/tools/readconfig"
)

func connect() (redis.Conn, error) {
	host := readconfig.ConfigData.String("redis::host")
	port := readconfig.ConfigData.String("redis::port")
	password := readconfig.ConfigData.String("redis::password")

	var Conn redis.Conn
	var err error
	if len(password) > 0 {
		Conn, err = redis.Dial("tcp", host+":"+port, redis.DialPassword(password))
	} else {
		Conn, err = redis.Dial("tcp", host+":"+port)
	}

	if err != nil {
		fmt.Println("redis connect error: " + fmt.Sprintf("%v", err))
	}

	return Conn, err
}

//获取key值
func Get(key string) (string, error) {
	rs, err := connect()
	if err != nil {
		return "", err
	}
	defer rs.Close()
	reply, err := rs.Do("GET", key)
	if reply == nil {
		return "", nil
	} else {
		return redis.String(reply, err)
	}
}

//设置key值
func Set(key string, value string) (interface{}, error) {
	rs, err := connect()
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	return rs.Do("SET", key, value)
}

//设置key值加过期时间
func SetWithSurvivalTime(key string, value string, existTime int) (interface{}, error) {
	if existTime == 0 {
		return Set(key, value)
	}

	rs, err := connect()
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	return rs.Do("SET", key, value, "EX", existTime)
}

//删除key值
func Del(key string) (interface{}, error) {
	rs, err := connect()
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	return rs.Do("DEL", key)
}

//向集合里添加元素
func SetAdd(key, value string) (interface{}, error) {
	rs, err := connect()
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	return rs.Do("SADD", key, value)
}

//返回集合里的元素列表
func SMEMBERS(key string) ([]string, error) {
	rs, err := connect()
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	reply, err := rs.Do("SMEMBERS", key)
	if reply == nil {
		return []string{}, nil
	} else {
		return redis.Strings(reply, err)
	}
}
