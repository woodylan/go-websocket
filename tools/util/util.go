package util

import (
	uuid "github.com/satori/go.uuid"
	"go-websocket/tools/readconfig"
	"strings"
)

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc, _ := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

//是否集群
func IsCluster() bool {
	cluster, _ := readconfig.ConfigData.Bool("common::cluster")

	return cluster
}