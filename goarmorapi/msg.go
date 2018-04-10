package goarmorapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type APIMsgRequest struct {
	Payload interface{} `json:"payload,omitempty"`
	Time    uint64      `json:"time,omitempty"`
}

type APIMsgResponse APIMsg

type APIMsg struct {
	Success bool        `json:"success"`
	Errs    []*APIErr   `json:"messages,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	Time    uint64      `json:"time,omitempty"`
}

type APIErrs []*APIErr

type APIErr struct {
	Code     uint64 `json:"code"`
	Err      error  `json:"message,omitempty"`
	Public   bool   `json:"-"`
	Severity uint64 `json:"-"`
}

type ErrSeverity uint64

const (
	ErrSeverityUnknown ErrSeverity = iota
	ErrSeverityDebug
	ErrSeverityInfo
	ErrSeverityWarn
	ErrSeverityError
	ErrSeverityFatal
	ErrSeverityPanic
)

func (v ErrSeverity) Uint64() uint64 {
	return uint64(v)
}

type APIErrCode uint64

func (v ErrSeverity) ErrorDefaultCode() uint64 {
	var u APIErrCode
	switch v {
	default:
		return 0
	case ErrSeverityDebug:
		u = ErrCodeDefaultDebug
	case ErrSeverityInfo:
		u = ErrCodeDefaultInfo
	case ErrSeverityWarn:
		u = ErrCodeDefaultWarn
	case ErrSeverityError:
		u = ErrCodeDefaultError
	case ErrSeverityFatal:
		u = ErrCodeDefaultFatal
	case ErrSeverityPanic:
		u = ErrCodeDefaultPanic
	}
	return uint64(u)
}

const (
	ErrCodeDefaultDebug APIErrCode = 1100
	ErrCodeDefaultInfo
	ErrCodeDefaultWarn

	ErrCodeDefaultError APIErrCode = 5100
	ErrCodeDefaultFatal
	ErrCodeDefaultPanic
)

func (v *APIMsg) Errors() []error {
	var a []error
	if v == nil || len(v.Errs) == 0 {
		return a
	}
	for _, e := range v.Errs {
		a = append(a, e)
	}
	return a
}

type ResponseMessenger interface {
	ResponseMessages() []*APIErr
}

func (v APIErrs) Errors() []error {
	var a []error
	for _, e := range v {
		a = append(a, e.Err)
	}
	return a
}

func (v APIErrs) First() error {
	a := []*APIErr(v)
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (v APIErrs) Last() error {
	a := []*APIErr(v)
	if len(a) == 0 {
		return nil
	}
	return a[len(a)-1]
}

func (e *APIErr) Error() string { return e.Err.Error() }

func (e *APIErr) MarshalJSON() ([]byte, error) {
	var (
		s string
		a []string
	)
	if e.Err != nil {
		s = e.Error()
		switch ErrSeverity(e.Severity) {
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

func (e *APIErr) UnmarshalJSON(b []byte) error {
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

func (a *APIMsg) KV() (KV, error) {
	if a == nil {
		return nil, errors.New("empty api response")
	}
	if len(a.Errs) == 0 {
		return nil, errors.New("empty key values")
	}
	kv := NewKV()
	for _, e := range a.Errs {
		if e.Code != uint64(ErrCodeDefaultDebug) {
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

func NewMsgRequest(
	responsePayload interface{}) (*APIMsgRequest, error) {
	return &APIMsgRequest{
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func NewMsgResponse(
	debugLevel int,
	isSuccess bool,
	responsePayload interface{},
	responseMessenger ResponseMessenger,
	errs ...*APIErr) (*APIMsgResponse, error) {
	publicErrs, err :=
		newMsgResponseErrs(debugLevel, responseMessenger, errs...)
	if err != nil {
		return nil, err
	}
	return &APIMsgResponse{
		Success: isSuccess,
		Errs:    publicErrs,
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func NewMsg(
	debugLevel int,
	isSuccess bool,
	responsePayload interface{},
	responseMessenger ResponseMessenger,
	errs ...*APIErr) (*APIMsg, error) {
	publicErrs, err :=
		newMsgResponseErrs(debugLevel, responseMessenger, errs...)
	if err != nil {
		return nil, err
	}
	return &APIMsg{
		Success: isSuccess,
		Errs:    publicErrs,
		Payload: responsePayload,
		Time:    uint64(time.Now().Unix())}, nil
}

func newMsgResponseErrs(
	debugLevel int,
	responseMessenger ResponseMessenger,
	errs ...*APIErr) ([]*APIErr, error) {
	errs = append(errs, responseMessenger.ResponseMessages()...)
	var publicErrs []*APIErr
	if debugLevel > 0 {
		for _, x := range errs {
			publicErrs = append(publicErrs,
				&APIErr{
					Code:     x.Code,
					Err:      errors.New(x.Error()),
					Public:   x.Public,
					Severity: x.Severity})
		}
	} else {
		isKVRemoved := false
		for _, x := range errs {
			if x.Public {
				publicErrs = append(publicErrs,
					&APIErr{
						Code:     x.Code,
						Err:      errors.New(x.Error()),
						Public:   x.Public,
						Severity: x.Severity})

				continue
			}
			if x.Code == uint64(ErrCodeDefaultDebug) {
				isKVRemoved = true
				continue
			}
			publicErrs = append(publicErrs,
				&APIErr{Code: x.Code, Severity: x.Severity})
		}
		if isKVRemoved {
			// Add empty (only with "code") "APIErr" structure in order to be able to
			// determine was an key-values in hadler's response.
			publicErrs = append(publicErrs, &APIErr{Code: uint64(ErrCodeDefaultDebug)})
		}
	}
	return publicErrs, nil
}
