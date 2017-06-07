package gameserverapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/armor5games/gameserver/gameserverconfigs"
	"github.com/armor5games/gameserver/gameservertokenverifiers"
	"github.com/sirupsen/logrus"
)

type SessionLoginPayload struct {
	UserID            uint64 `json:"userID,omitempty"`
	UserDataVersion   uint64 `json:"userDataVersion,omitempty"`
	AccessToken       string `json:"accessToken,omitempty"`
	ShardURL          string `json:"shardURL,omitempty"`
	StaticDataVersion uint64 `json:"staticDataVersion,omitempty"`
	UserName          string `json:"userName,omitempty"`
	NewUser           bool   `json:"newUser"`
}

type SessionShardResponse struct {
	JSON
	Payload *SessionShardPayload `json:",omitempty"`
}

type SessionShardPayload struct {
	AccessToken         string `json:"accessToken,omitempty"`
	AccessTokenVerifier string `json:"accessTokenVerifier,omitempty"`
	UserDataVersion     uint64 `json:"userDataVersion,omitempty"`
}

func (ssr *SessionShardPayload) IsAccessTokenPresent() bool {
	return len(ssr.AccessToken) > 16
}

func (ssr *SessionShardPayload) IsUserDataVersionPresent() bool {
	return ssr.UserDataVersion != 0
}

func NewSession(ctx context.Context, shardURL *url.URL) (
	*SessionShardResponse, error) {
	config, ok := ctx.Value(CtxConfigKey).(*gameserverconfigs.Config)
	if !ok {
		return nil, errors.New("context.Value fn error")
	}

	netClient := &http.Client{
		Timeout: time.Second * time.Duration(config.Server.APITimeoutSeconds)}
	res, err := netClient.Post(shardURL.String(),
		"application/x-www-form-urlencoded", strings.NewReader(shardURL.RawQuery))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	httpBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll fn error: %s", err.Error())
	}

	shardResponse := &SessionShardResponse{}
	err = json.Unmarshal(httpBody, shardResponse)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal fn error: %s", err.Error())
	}

	l, ok := ctx.Value(CtxLoggerWithoutUserKey).(*logrus.Logger)
	if !ok {
		return nil, errors.New("context.Value fn error")
	}

	kv := NewKV()
	kv["httpStatusCode"] = res.StatusCode
	kv["sessionSuccessStatus"] = shardResponse.Success

	if res.StatusCode != 200 || !shardResponse.Success {
		err = errors.New("session status error")
		l.WithFields(logrus.Fields(kv)).Error(err.Error())
	}

	if shardResponse.Payload == nil {
		err = errors.New("empty session payload")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	userNewDataVersion := shardResponse.Payload.UserDataVersion
	kv["userNewDataVersion"] = userNewDataVersion

	if !shardResponse.Payload.IsUserDataVersionPresent() {
		err = errors.New("shard server return zero user data version")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	accessToken := shardResponse.Payload.AccessToken
	kv["accessToken"] = accessToken

	if !shardResponse.Payload.IsAccessTokenPresent() {
		err = errors.New("shard server return zero access token")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	accessTokenVerifier :=
		gameservertokenverifiers.New(accessToken, config.Server.ServerSecretKey)
	kv["accessTokenVerifier"] = accessTokenVerifier

	if shardResponse.Payload.AccessTokenVerifier != accessTokenVerifier {
		err = errors.New("AccessToken and AccessTokenVerifier mismatch")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	return shardResponse, nil
}
