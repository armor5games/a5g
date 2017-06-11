package goarmorconfigs

import (
	"crypto/md5"
	"errors"
	"fmt"
)

func (c *Config) NewDummyChecksum(s string) (string, error) {
	if c == nil || c.Server.ServerSecretKey == "" {
		return "", errors.New("empty secret")
	}

	if s == "" {
		return "", errors.New("empty token")
	}

	return fmt.Sprintf("%x", md5.Sum(
		[]byte(fmt.Sprintf("%s%s", s, c.Server.ServerSecretKey)))), nil
}
