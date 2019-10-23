package src

import (
	uuid "github.com/satori/go.uuid"
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
