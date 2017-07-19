package goarmorvkpayments

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JSONResponse struct {
	Error   *JSONError   `json:"error,omitempty"`
	Payload *JSONPayload `json:"response,omitempty"`
}

// JSONError <https://vk.com/dev/payments_errors>.
type JSONError struct {
	Code     uint64 `json:"error_code"`
	Err      error  `json:"error_msg,omitempty"`
	Critical bool   `json:"critical"`
}

// JSONPayload <https://vk.com/dev/payments_getitem>,
// <https://vk.com/dev/payments_status>.
type JSONPayload struct {
	Title    string `json:"title,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
	Price    uint64 `json:"price,omitempty"`
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("%d critical=%t %s", e.Code, e.Critical, e.Err.Error())
}

func (e *JSONError) MarshalJSON() ([]byte, error) {
	var m string

	if e.Err != nil {
		m = e.Err.Error()
	}

	return json.Marshal(&struct {
		Code     uint64 `json:"error_code"`
		Message  string `json:"error_msg,omitempty"`
		Critical bool   `json:"critical"`
	}{
		Code:     e.Code,
		Message:  m,
		Critical: e.Critical})
}

func (e *JSONError) UnmarshalJSON(b []byte) error {
	s := &struct {
		Code     uint64 `json:"error_code"`
		Message  string `json:"error_msg,omitempty"`
		Critical bool   `json:"critical"`
	}{}

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	e.Code = s.Code
	e.Critical = s.Critical

	if s.Message != "" {
		e.Err = errors.New(s.Message)
	}

	return nil
}
