package goarmorchecksums

import (
	"github.com/pkg/errors"
)

var (
	ErrPayloadEmpty   = errors.New("empty payload")
	ErrSecretKeyEmpty = errors.New("empty secret key")
)
