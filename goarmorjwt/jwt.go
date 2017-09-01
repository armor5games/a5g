package goarmorjwt

import (
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var ErrSecretKeyEmpty = errors.New("empty secret key")

func NewSession(userID int64, tokenDuration time.Duration, secretKey string) (
	string, *jwt.StandardClaims, error) {
	if secretKey == "" {
		return "", nil, errors.WithStack(ErrSecretKeyEmpty)
	}

	t := time.Now()

	sessionClaims := &jwt.StandardClaims{
		ExpiresAt: t.Add(tokenDuration).Unix(),
		IssuedAt:  t.Unix(),
		Issuer:    strconv.FormatInt(userID, 10),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims)
	accessToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, errors.Wrap(err, "jwt.(*Token).SignedString fn")
	}

	return accessToken, sessionClaims, nil
}
