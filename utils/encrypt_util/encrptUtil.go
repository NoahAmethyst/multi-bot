package encrypt_util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"strings"
)

func GenerateUuid(format bool) string {
	uuidValue := uuid.New().String()
	if format {
		uuidValue = strings.Replace(uuidValue, "-", "", -1)
	}
	return uuidValue
}

func Md5WithSecretKey(content string, secretKey string) string {
	v := content + secretKey
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
