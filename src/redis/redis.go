package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-websocket/src/readConfig"
)

func connect() (redis.Conn, error) {
	host := readConfig.ConfigData.String("redis::host")
	port := readConfig.ConfigData.String("redis::port")
	password := readConfig.ConfigData.String("redis::password")

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
func SetMembers(key string) ([]string, error) {
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
