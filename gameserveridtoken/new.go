package gameserveridtoken

import (
	"crypto/md5"
	"fmt"
)

func New(id, idSecret string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", id, idSecret))))
}
