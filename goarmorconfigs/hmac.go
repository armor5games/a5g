package goarmorconfigs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

const HMACAlgorithm = "sha256"

func (c *Config) HMACChecksum(payload []byte) (string, error) {
	serverSecretKey := c.Server.ServerSecretKey
	if serverSecretKey == "" {
		return "", errors.New("empty ServerSecretKey config")
	}

	hmacChecksum := hmac.New(sha256.New, []byte(serverSecretKey))

	_, err := hmacChecksum.Write(payload)
	if err != nil {
		return "", fmt.Errorf("hash.Hash.Write fn: %s", err.Error())
	}

	return hex.EncodeToString(hmacChecksum.Sum(nil)), nil
}
