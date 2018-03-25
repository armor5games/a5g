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
	Success  bool        `json:"success"`
	Messages []*JSONMsg  `json:"messages,omitempty"`
	Payload  interface{} `json:"payload,omitempty"`
	Time     uint64      `json:"time,omitempty"`
}

type JSONMessages []*JSONMsg

type JSONMsg struct {
	Code     uint64 `json:"code"`
	Err      error  `json:"message,omitempty"`
	Public   bool   `json:"-"`
	Severity uint64 `json:"-"`
}

type JSONMsgSeverity uint64

const (
	MsgSeverityUnknown JSONMsgSeverity = iota
	MsgSeverityDebug
	MsgSeverityInfo
	MsgSeverityWarn
	MsgSeverityError
	MsgSeverityFatal
	MsgSeverityPanic
)

func (v JSONMsgSeverity) Uint64() uint64 {
	return uint64(v)
}

func (v JSONMsgSeverity) ErrorDefaultCode() uint64 {
	var u JSONMsgCode
	switch v {
	default:
		return 0
	case MsgSeverityDebug:
		u = MsgCodeDefaultDebug
	case MsgSeverityInfo:
		u = MsgCodeDefaultInfo
	case MsgSeverityWarn:
		u = MsgCodeDefaultWarn
	case MsgSeverityError:
		u = MsgCodeDefaultError
	case MsgSeverityFatal:
		u = MsgCodeDefaultFatal
	case MsgSeverityPanic:
		u = MsgCodeDefaultPanic
	}
	return uint64(u)
}

type JSONMsgCode uint64

const (
	MsgCodeDefaultDebug JSONMsgCode = 1100
	MsgCodeDefaultInfo
	MsgCodeDefaultWarn

	MsgCodeDefaultError JSONMsgCode = 5100
	MsgCodeDefaultFatal
	MsgCodeDefaultPanic
)

func (v *JSONResponse) Errors() []error {
	var a []error
	if v == nil || len(v.Messages) == 0 {
		return a
	}
	for _, e := range v.Messages {
		a = append(a, e)
	}
	return a
}

type ResponseMessenger interface {
	ResponseMessages() []*JSONMsg
}

func (v JSONMessages) Errors() []error {
	var a []error
	for _, e := range v {
		a = append(a, e.Err)
	}
	return a
}

func (v JSONMessages) First() error {
	a := []*JSONMsg(v)
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (v JSONMessages) Last() error {
	a := []*JSONMsg(v)
	if len(a) == 0 {
		return nil
	}
	return a[len(a)-1]
}

func (e *JSONMsg) Error() string { return e.Err.Error() }

func (e *JSONMsg) MarshalJSON() ([]byte, error) {
	var (
		s string
		a []string
	)
	if e.Err != nil {
		s = e.Error()
		switch JSONMsgSeverity(e.Severity) {
		case MsgSeverityError, MsgSeverityFatal, MsgSeverityPanic:
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

func (e *JSONMsg) UnmarshalJSON(b []byte) error {
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

func (a *JSONResponse) KV() (KV, error) {
	if a == nil {
		return nil, errors.New("empty api response")
	}
	if len(a.Messages) == 0 {
		return nil, errors.New("empty key values")
	}
	kv := NewKV()
	for _, e := range a.Messages {
		if e.Code != uint64(MsgCodeDefaultDebug) {
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
	responseMessenger ResponseMessenger,
	errs ...*JSONMsg) (*JSONResponse, error) {
	publicMessages, err :=
		newJSONResponseMessages(debugLevel, responseMessenger, errs...)
	if err != nil {
		return nil, err
	}
	return &JSONResponse{
		Success:  isSuccess,
		Messages: publicMessages,
		Payload:  responsePayload,
		Time:     uint64(time.Now().Unix())}, nil
}

func newJSONResponseMessages(
	debugLevel int,
	responseMessenger ResponseMessenger,
	errs ...*JSONMsg) ([]*JSONMsg, error) {
	errs = append(errs, responseMessenger.ResponseMessages()...)
	var publicMessages []*JSONMsg
	if debugLevel > 0 {
		for _, x := range errs {
			publicMessages = append(publicMessages,
				&JSONMsg{
					Code:     x.Code,
					Err:      errors.New(x.Error()),
					Public:   x.Public,
					Severity: x.Severity})
		}
	} else {
		isKVRemoved := false
		for _, x := range errs {
			if x.Public {
				publicMessages = append(publicMessages,
					&JSONMsg{
						Code:     x.Code,
						Err:      errors.New(x.Error()),
						Public:   x.Public,
						Severity: x.Severity})

				continue
			}
			if x.Code == uint64(MsgCodeDefaultDebug) {
				isKVRemoved = true
				continue
			}
			publicMessages = append(publicMessages,
				&JSONMsg{Code: x.Code, Severity: x.Severity})
		}
		if isKVRemoved {
			// Add empty (only with "code") "JSONMsg" structure in order to be able to
			// determine was an key-values in hadler's response.
			publicMessages = append(publicMessages, &JSONMsg{Code: uint64(MsgCodeDefaultDebug)})
		}
	}
	return publicMessages, nil
}
