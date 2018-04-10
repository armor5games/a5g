package goarmorapi

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SessionShardResponse struct {
	APIMsg
	Payload *SessionShardPayload `json:"payload,omitempty"`
}

type SessionShardPayload struct {
	UserDataVersion uint64 `json:"userDataVersion,omitempty"`
}

func (v *SessionShardPayload) IsUserDataVersionPresent() bool {
	return v.UserDataVersion != 0
}

func NewSession(
	appLogger *logrus.Logger,
	requestToShard *http.Request,
	timeoutDuration time.Duration,
	secretKey string) (*SessionShardResponse, error) {
	if secretKey == "" {
		return nil, errors.New("empty secret key")
	}

	c := &http.Client{Timeout: timeoutDuration}

	res, err := c.Do(requestToShard)
	if err != nil {
		return nil, err
	}
	defer func(l *logrus.Logger) {
		if err := res.Body.Close(); err != nil {
			l.Error(err.Error())
		}
	}(appLogger)

	httpBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll fn")
	}

	kv := NewKV()
	kv["httpStatusCode"] = res.StatusCode

	if res.StatusCode != http.StatusOK {
		err = errors.New("session http status code error")
		appLogger.WithFields(logrus.Fields(kv)).Error(err.Error())
	}

	kv["httpContentType"] = res.Header.Get("Content-type")

	t, _, err := mime.ParseMediaType(res.Header.Get("Content-type"))
	if err != nil {
		return nil, errors.Wrap(err, "mime.ParseMediaType fn")
	}

	if t != "application/json" {
		err = errors.New("shard server's response is not in json format")

		kv["httpBody"] = string(httpBody)
		appLogger.WithFields(logrus.Fields(kv)).Error(err.Error())

		return nil, err
	}

	var shardResponse = new(SessionShardResponse)

	err = json.Unmarshal(httpBody, shardResponse)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal fn")
	}

	if !shardResponse.Success {
		err = errors.New("session is not successful")
		appLogger.WithFields(logrus.Fields(kv)).Error(err.Error())
	}

	kv["sessionSuccessStatus"] = shardResponse.Success

	if shardResponse.Payload == nil {
		err = errors.New("empty session payload")

		appLogger.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	userNewDataVersion := shardResponse.Payload.UserDataVersion
	kv["userNewDataVersion"] = userNewDataVersion

	if !shardResponse.Payload.IsUserDataVersionPresent() {
		err = errors.New("shard server retruns zero user's data version")

		appLogger.WithFields(logrus.Fields(kv)).Error(err.Error())

		return shardResponse, err
	}

	return shardResponse, nil
}
