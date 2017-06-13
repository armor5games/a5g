package goarmorapi

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/armor5games/goarmor/goarmorconfigs"
)

type JSON struct {
	Success bool
	Errors  []*ErrorJSON `json:",omitempty"`
	Payload interface{}  `json:",omitempty"`
	Time    uint64       `json:",omitempty"`
}

type ErrorJSON struct {
	Code   uint64
	Error  error `json:"Message,omitempty"`
	Public bool  `json:"-"`
}

func (e *ErrorJSON) MarshalJSON() ([]byte, error) {
	var m string

	if e.Error != nil {
		m = e.Error.Error()
	}

	return json.Marshal(&struct {
		Code    uint64
		Message string `json:",omitempty"`
	}{
		Code:    e.Code,
		Message: m})
}

func (e *ErrorJSON) UnmarshalJSON(b []byte) error {
	s := &struct {
		Code    uint64
		Message string
	}{}

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	e.Code = s.Code

	if s.Message != "" {
		e.Error = errors.New(s.Message)
	}

	return nil
}

func (j *JSON) KV() (KV, error) {
	if j == nil {
		return nil, errors.New("empty api response")
	}

	if len(j.Errors) == 0 {
		return nil, errors.New("empty key values")
	}

	kv := NewKV()

	for _, e := range j.Errors {
		if e.Code != KVAPIErrorCode {
			continue
		}

		if e.Error.Error() == "" {
			return nil, errors.New("empty kv")
		}

		x := strings.SplitN(e.Error.Error(), ":", 2)
		if len(x) != 2 {
			return nil, errors.New("bad kv format")
		}

		kv[x[0]] = x[1]
	}

	if len(kv) == 0 {
		return nil, errors.New("empty kv")
	}

	return kv, nil
}

func ResponseJSON(
	ctx context.Context,
	isSuccess bool,
	responsePayload interface{},
	keyValues KV,
	errs ...*ErrorJSON) (*JSON, error) {
	config, ok := ctx.Value(CtxKeyConfig).(*goarmorconfigs.Config)
	if !ok {
		return nil, errors.New("context.Value fn error")
	}

	if len(keyValues) != 0 {
		errs = append(errs, keyValues.ResponseErrors()...)
	}

	publicErrors := make([]*ErrorJSON, 0)

	if config.Server.DebuggingLevel > 0 {
		publicErrors = errs
	} else {
		isKVRemoved := false

		for _, x := range errs {
			if x.Public == true {
				publicErrors = append(publicErrors, x)

				continue
			}

			if x.Code == 1100 {
				isKVRemoved = true

				continue
			}

			x.Error = nil
			publicErrors = append(publicErrors, x)
		}

		if isKVRemoved {
			// Add empty (only with "code") "ErrorJSON" structure in order to be able to
			// determine was an key-values in hadler's response.
			publicErrors = append(publicErrors, &ErrorJSON{Code: 1100})
		}
	}

	return &JSON{
		Success: isSuccess,
		Errors:  publicErrors,
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}
