package goarmorconfigs

import (
	"crypto/md5"
	"errors"
	"fmt"
)

func (c *Config) NewDummyChecksum(s string) (string, error) {
	if c == nil || c.Server.ServerSecretKey == "" {
		return "", errors.New("Empty secret")
	}

	if s == "" {
		return "", errors.New("Empty token")
	}

	return fmt.Sprintf("%x", md5.Sum(
		[]byte(fmt.Sprintf("%s%s", s, c.Server.ServerSecretKey)))), nil
}

func (c *Config) NewDummyUserChecksum(payload []byte, secure []byte) (
	string, error) {
	if c == nil || c.ShardServer.USRSec == "" {
		return "", errors.New("Empty user secret")
	}

	return fmt.Sprintf("%x", md5.Sum(
		[]byte(fmt.Sprintf("%s%s%s", payload, secure, c.ShardServer.USRSec)))), nil
}
