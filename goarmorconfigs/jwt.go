package goarmorconfigs

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func (c *Config) NewSession(userID int64, tokenDuration time.Duration) (
	string, *jwt.StandardClaims, error) {
	t := time.Now()

	sessionClaims := &jwt.StandardClaims{
		ExpiresAt: t.Add(tokenDuration).Unix(),
		IssuedAt:  t.Unix(),
		Issuer:    strconv.FormatInt(userID, 10),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims)
	accessToken, err := jwtToken.
		SignedString([]byte(c.Server.ServerSecretKey))
	if err != nil {
		return "", nil, fmt.Errorf("jwt.(*Token).SignedString fn error: %s",
			err.Error())
	}

	return accessToken, sessionClaims, nil
}
