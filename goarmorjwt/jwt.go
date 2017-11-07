package goarmorjwt

import (
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var ErrSecretKeyEmpty = errors.New("empty secret key")

func NewSession(
	secretKey string, userID int64, issuedAt time.Time, lifeTime time.Duration) (
	string, *jwt.StandardClaims, error) {
	if secretKey == "" {
		return "", nil, errors.WithStack(ErrSecretKeyEmpty)
	}

	sessionClaims := &jwt.StandardClaims{
		ExpiresAt: issuedAt.Add(lifeTime).Unix(),
		IssuedAt:  issuedAt.Unix(),
		Issuer:    strconv.FormatInt(userID, 10),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims)
	accessToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, errors.Wrap(err, "jwt.(*Token).SignedString fn")
	}

	return accessToken, sessionClaims, nil
}
