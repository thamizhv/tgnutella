package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func GetHash(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}
