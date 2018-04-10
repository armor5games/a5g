package goarmorchecksums

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
)

const HMACAlgorithm = "sha256"

func NewHMAC(toCheck []byte, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.WithStack(ErrSecretKeyEmpty)
	}
	hmacChecksum := hmac.New(sha256.New, []byte(secretKey))
	_, err := hmacChecksum.Write(toCheck)
	if err != nil {
		return "", errors.Wrap(err, "hash.Hash.Write fn")
	}
	return hex.EncodeToString(hmacChecksum.Sum(nil)), nil
}
