package gameserverapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/armor5games/gameserver/gameserverconfigs"
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
	AccessTokenChecksum string `json:"accessTokenChecksum,omitempty"`
	UserDataVersion     uint64 `json:"userDataVersion,omitempty"`
}

func (ssr *SessionShardPayload) IsAccessTokenPresent() bool {
	return len(ssr.AccessToken) >= 16
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

	l, ok := ctx.Value(CtxLoggerKey).(*logrus.Logger)
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
	defer func(l *logrus.Logger) {
		if err := res.Body.Close(); err != nil {
			l.Error(err.Error())
		}
	}(l)

	httpBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll fn error: %s", err.Error())
	}

	kv := NewKV()
	kv["httpStatusCode"] = res.StatusCode

	if res.StatusCode != http.StatusOK {
		err = errors.New("session http status code error")
		l.WithFields(logrus.Fields(kv)).Error(err.Error())
	}

	kv["httpContentType"] = res.Header.Get("Content-type")

	t, _, err := mime.ParseMediaType(res.Header.Get("Content-type"))
	if err != nil {
		return nil, fmt.Errorf("mime.ParseMediaType fn error: %s", err.Error())
	}

	if t != "application/json" {
		err = errors.New("shard server's response is not in json format")

		kv["httpBody"] = string(httpBody)
		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return nil, err
	}

	shardResponse := &SessionShardResponse{}
	err = json.Unmarshal(httpBody, shardResponse)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal fn error: %s", err.Error())
	}

	if !shardResponse.Success {
		err = errors.New("session is not successful")
		l.WithFields(logrus.Fields(kv)).Error(err.Error())
	}

	kv["sessionSuccessStatus"] = shardResponse.Success

	if shardResponse.Payload == nil {
		err = errors.New("empty session payload")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	userNewDataVersion := shardResponse.Payload.UserDataVersion
	kv["userNewDataVersion"] = userNewDataVersion

	if !shardResponse.Payload.IsUserDataVersionPresent() {
		err = errors.New("shard server return zero user's data version")

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

	accessTokenChecksum, err := config.NewDummyChecksum(accessToken)
	if shardResponse.Payload.AccessTokenChecksum != accessTokenChecksum {
		err = fmt.Errorf("gameserverconfigs.(*Config)NewDummyChecksum fn: %s",
			err.Error())

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	kv["accessTokenChecksum"] = accessTokenChecksum

	if shardResponse.Payload.AccessTokenChecksum != accessTokenChecksum {
		err = errors.New("AccessToken and AccessTokenChecksum mismatch")

		l.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	return shardResponse, nil
}
