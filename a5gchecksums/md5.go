package a5gchecksums

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"

	"github.com/pkg/errors"
)

func NewMD5(toCheck []byte, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.WithStack(ErrSecretKeyEmpty)
	}

	buf := bytes.NewBuffer(toCheck)

	_, err := buf.WriteString(secretKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	a := md5.Sum(buf.Bytes())

	if len(toCheck) == 0 {
		return hex.EncodeToString(a[:]), errors.WithStack(ErrPayloadEmpty)
	}

	return hex.EncodeToString(a[:]), nil
}

func NewMD5Salted(toCheck []byte, secretKey, checksumSalt string) (
	string, error) {
	if secretKey == "" {
		return "", errors.WithStack(ErrSecretKeyEmpty)
	}

	if checksumSalt == "" {
		return "", errors.New("missing salt")
	}

	buf := bytes.NewBuffer(toCheck)

	_, err := buf.WriteString(checksumSalt)
	if err != nil {
		return "", errors.WithStack(err)
	}

	_, err = buf.WriteString(secretKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	a := md5.Sum(buf.Bytes())

	if len(toCheck) == 0 {
		return hex.EncodeToString(a[:]), errors.WithStack(ErrPayloadEmpty)
	}

	return hex.EncodeToString(a[:]), nil
}
