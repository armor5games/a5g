package goarmorapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type JSONRequest struct {
	Payload interface{} `json:"payload,omitempty"`
	Time    uint64      `json:"time,omitempty"`
}

type JSONResponse struct {
	Success bool         `json:"success"`
	Errors  []*ErrorJSON `json:"messages,omitempty"`
	Payload interface{}  `json:"payload,omitempty"`
	Time    uint64       `json:"time,omitempty"`
}

type ErrorsJSON []*ErrorJSON

type ErrorJSON struct {
	Code     uint64 `json:"code"`
	Err      error  `json:"message,omitempty"`
	Public   bool   `json:"-"`
	Severity uint64 `json:"-"`
}

type ErrorJSONSeverity uint64

const (
	ErrSeverityUnknown ErrorJSONSeverity = iota
	ErrSeverityDebug
	ErrSeverityInfo
	ErrSeverityWarn
	ErrSeverityError
	ErrSeverityFatal
	ErrSeverityPanic
)

func (v ErrorJSONSeverity) Uint64() uint64 {
	return uint64(v)
}

func (v ErrorJSONSeverity) ErrorDefaultCode() uint64 {
	var u ErrorJSONCode

	switch v {
	default:
		return 0

	case ErrSeverityDebug:
		u = ErrCodeDefautlDebug

	case ErrSeverityInfo:
		u = ErrCodeDefautlInfo

	case ErrSeverityWarn:
		u = ErrCodeDefautlWarn

	case ErrSeverityError:
		u = ErrCodeDefautlError

	case ErrSeverityFatal:
		u = ErrCodeDefautlFatal

	case ErrSeverityPanic:
		u = ErrCodeDefautlPanic
	}

	return uint64(u)
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

func (errorsJSON ErrorsJSON) First() error {
	a := []*ErrorJSON(errorsJSON)

	if len(a) == 0 {
		return nil
	}

	return a[0]
}

func (errorsJSON ErrorsJSON) Last() error {
	a := []*ErrorJSON(errorsJSON)

	if len(a) == 0 {
		return nil
	}

	return a[len(a)-1]
}

func (e *ErrorJSON) Error() string {
	return e.Err.Error()
}

func (e *ErrorJSON) MarshalJSON() ([]byte, error) {
	var (
		s string
		a []string
	)

	if e.Err != nil {
		s = e.Error()

		switch ErrorJSONSeverity(e.Severity) {
		case ErrSeverityError, ErrSeverityFatal, ErrSeverityPanic:
			a = strings.Split(fmt.Sprintf("%+v", e.Err), "\n")
			a = append(a[1:2], a[2:]...)
		}
	}

	return json.Marshal(&struct {
		Code       uint64   `json:"code,omitempty"`
		Message    string   `json:"message,omitempty"`
		StackTrace []string `json:"stackTrace,omitempty"`
	}{
		Code:       e.Code,
		Message:    s,
		StackTrace: a})
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
		e.Err = errors.New(s.Message)
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

		if e.Error() == "" {
			return nil, errors.New("empty kv")
		}

		x := strings.SplitN(e.Error(), ":", 2)
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
	responsePayload interface{}) (*JSONRequest, error) {
	return &JSONRequest{
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func NewJSONResponse(
	debugLevel int,
	isSuccess bool,
	responsePayload interface{},
	responseErrorer ResponseErrorer,
	errs ...*ErrorJSON) (*JSONResponse, error) {
	publicErrors, err :=
		newJSONResponseErrors(debugLevel, responseErrorer, errs...)
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
	debugLevel int,
	responseErrorer ResponseErrorer,
	errs ...*ErrorJSON) ([]*ErrorJSON, error) {
	errs = append(errs, responseErrorer.ResponseErrors()...)

	var publicErrors []*ErrorJSON

	if debugLevel > 0 {
		for _, x := range errs {
			publicErrors = append(publicErrors,
				&ErrorJSON{
					Code:     x.Code,
					Err:      errors.New(x.Error()),
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
						Err:      errors.New(x.Error()),
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
