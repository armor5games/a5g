package goarmorchecksums

import (
	"bytes"
	"crypto/md5"

	"github.com/pkg/errors"
)

var (
	ErrPayloadEmpty   = errors.New("empty payload")
	ErrSecretKeyEmpty = errors.New("empty secret key")
)

func New(toCheck []byte, secretKey string) ([]byte, error) {
	if secretKey == "" {
		return nil, errors.WithStack(ErrSecretKeyEmpty)
	}

	buf := bytes.NewBuffer(toCheck)

	_, err := buf.WriteString(secretKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	x := md5.Sum(buf.Bytes())

	if len(toCheck) == 0 {
		return x[:], errors.WithStack(ErrPayloadEmpty)
	}

	return x[:], nil
}

func NewWithSalt(toCheck []byte, secretKey, checksumSalt string) (
	[]byte, error) {
	if secretKey == "" {
		return nil, errors.WithStack(ErrSecretKeyEmpty)
	}

	if checksumSalt == "" {
		return nil, errors.New("missing salt")
	}

	buf := bytes.NewBuffer(toCheck)

	_, err := buf.WriteString(checksumSalt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = buf.WriteString(secretKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	x := md5.Sum(buf.Bytes())

	if len(toCheck) == 0 {
		return x[:], errors.WithStack(ErrPayloadEmpty)
	}

	return x[:], nil
}
