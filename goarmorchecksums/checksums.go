package goarmorchecksums

import (
	"crypto/md5"
	"errors"
	"fmt"
)

func New(payload []byte, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.New("Empty secret key")
	}

	if len(payload) == 0 {
		return "", errors.New("Empty token")
	}

	return fmt.Sprintf("%x", md5.Sum(
		[]byte(fmt.Sprintf("%s%s", payload, secretKey)))), nil
}

func NewWithSalt(payload []byte, secretKey, secretSalt string) (
	string, error) {
	if secretKey == "" {
		return "", errors.New("Empty secret key")
	}

	if secretSalt == "" {
		return "", errors.New("Empty salt")
	}

	s := fmt.Sprintf("%s%s%s", payload, secretSalt, secretKey)

	return fmt.Sprintf("%x", md5.Sum([]byte(s))), nil
}
