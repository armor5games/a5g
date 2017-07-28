package goarmorapi

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/armor5games/goarmor/goarmorconfigs"
)

type JSONRequest struct {
	Payload interface{} `json:"payload,omitempty"`
	Time    uint64      `json:"time,omitempty"`
}

type JSONResponse struct {
	Success bool         `json:"success,omitempty"`
	Errors  []*ErrorJSON `json:"messages,omitempty"`
	Payload interface{}  `json:"payload,omitempty"`
	Time    uint64       `json:"time,omitempty"`
}

type ErrorJSON struct {
	Code uint64
	// TODO: Rename "Error" to "Err"
	Error    error  `json:"message,omitempty"`
	Public   bool   `json:"-"`
	Severity uint64 `json:"-"`
}

type ErrorJSONSeverity uint64

const (
	ErrSeverityDebug ErrorJSONSeverity = iota
	ErrSeverityInfo
	ErrSeverityWarn
	ErrSeverityError
	ErrSeverityFatal
	ErrSeverityPanic
)

func (s ErrorJSONSeverity) ErrorDefaultCode() uint64 {
	var c ErrorJSONCode

	switch s {
	default:
		return 0

	case ErrSeverityDebug:
		c = ErrCodeDefautlDebug

	case ErrSeverityInfo:
		c = ErrCodeDefautlInfo

	case ErrSeverityWarn:
		c = ErrCodeDefautlWarn

	case ErrSeverityError:
		c = ErrCodeDefautlError

	case ErrSeverityFatal:
		c = ErrCodeDefautlFatal

	case ErrSeverityPanic:
		c = ErrCodeDefautlPanic
	}

	return uint64(c)
}

type ErrorJSONCode uint64

const (
	ErrCodeDefautlDebug ErrorJSONCode = 1100
	ErrCodeDefautlInfo
	ErrCodeDefautlWarn

	ErrCodeDefautlError ErrorJSONCode = 5100
	ErrCodeDefautlFatal
	ErrCodeDefautlPanic
)

type ResponseErrorer interface {
	ResponseErrors() []*ErrorJSON
}

func (e *ErrorJSON) MarshalJSON() ([]byte, error) {
	var m string

	if e.Error != nil {
		m = e.Error.Error()
	}

	return json.Marshal(&struct {
		Code    uint64 `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}{
		Code:    uint64(e.Code),
		Message: m})
}

func (e *ErrorJSON) UnmarshalJSON(b []byte) error {
	s := &struct {
		Code    uint64 `json:"code"`
		Message string `json:"message"`
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

func (j *JSONResponse) KV() (KV, error) {
	if j == nil {
		return nil, errors.New("empty api response")
	}

	if len(j.Errors) == 0 {
		return nil, errors.New("empty key values")
	}

	kv := NewKV()

	for _, e := range j.Errors {
		if e.Code != uint64(ErrCodeDefautlDebug) {
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

func NewJSONRequest(
	ctx context.Context,
	responsePayload interface{}) (*JSONRequest, error) {
	return &JSONRequest{
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func NewJSONResponse(
	ctx context.Context,
	isSuccess bool,
	responsePayload interface{},
	responseErrorer ResponseErrorer,
	errs ...*ErrorJSON) (*JSONResponse, error) {
	publicErrors, err := newJSONResponseErrors(ctx, responseErrorer, errs...)
	if err != nil {
		return nil, err
	}

	return &JSONResponse{
		Success: isSuccess,
		Errors:  publicErrors,
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func newJSONResponseErrors(
	ctx context.Context,
	responseErrorer ResponseErrorer,
	errs ...*ErrorJSON) ([]*ErrorJSON, error) {
	config, ok := ctx.Value(CtxKeyConfig).(goarmorconfigs.Configer)
	if !ok {
		return nil, errors.New("context.Value fn error")
	}

	errs = append(errs, responseErrorer.ResponseErrors()...)

	var publicErrors []*ErrorJSON

	if config.ServerDebuggingLevel() > 0 {
		for _, x := range errs {
			publicErrors = append(publicErrors,
				&ErrorJSON{
					Code:     x.Code,
					Error:    errors.New(x.Error.Error()),
					Public:   x.Public,
					Severity: x.Severity})
		}

	} else {
		isKVRemoved := false

		for _, x := range errs {
			if x.Public {
				publicErrors = append(publicErrors,
					&ErrorJSON{
						Code:     x.Code,
						Error:    errors.New(x.Error.Error()),
						Public:   x.Public,
						Severity: x.Severity})

				continue
			}

			if x.Code == uint64(ErrCodeDefautlDebug) {
				isKVRemoved = true

				continue
			}

			publicErrors = append(publicErrors,
				&ErrorJSON{Code: x.Code, Severity: x.Severity})
		}

		if isKVRemoved {
			// Add empty (only with "code") "ErrorJSON" structure in order to be able to
			// determine was an key-values in hadler's response.
			publicErrors = append(publicErrors, &ErrorJSON{Code: uint64(ErrCodeDefautlDebug)})
		}
	}

	return publicErrors, nil
}
