package gameservertokenverifiers

import (
	"crypto/md5"
	"fmt"
)

func New(tokenStr, tokenSecret string) string {
	return fmt.Sprintf("%x", md5.Sum(
		[]byte(fmt.Sprintf("%s%s", tokenStr, tokenSecret))))
}
