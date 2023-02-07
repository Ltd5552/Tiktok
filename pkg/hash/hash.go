package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5WithSalt(str string, salt string) string {
	h := md5.New()
	h.Write([]byte(salt))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
