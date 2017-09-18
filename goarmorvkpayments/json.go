package goarmorvkpayments

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JSONResponse struct {
	Error   *JSONError  `json:"Error,omitempty"`
	Payload interface{} `json:"Response,omitempty"`
}

// JSONError <https://vk.com/dev/payments_errors>.
type JSONError struct {
	Code     uint64 `json:"error_code"`
	Err      error  `json:"error_msg,omitempty"`
	Critical bool   `json:"critical"`
}

// JSONSuccessInfo <https://vk.com/dev/payments_getitem>.
type JSONSuccessInfo struct {
	// Title vk description: название товара, до 48 символов
	Title string `json:"title"`
	// PhotoURL vk description: URL изображения товара на сервере
	// разработчика. Рекомендуемый размер изображения – 75х75px.
	PhotoURL string `json:"photo_url,omitempty"`
	// Price vk description: стоимость товара в голосах.
	Price int64 `json:"price"`
	// ItemID vk description: идентификатор товара в приложении.
	ItemID string `json:"item_id,omitempty"`
	// Expiration vk description: разрешает кэширование товара на
	// {expiration} секунд. Допустимый диапазон от 600 до 604800 секунд.
	// Внимание! При отсутствии параметра возможно кэширование товара на
	// 3600 секунд при большом количестве подряд одинаковых ответов. Для
	// отмены кэширования необходимо передать 0 в качестве значения
	// параметра.
	Expiration int64 `json:"expiration"`
}

type JSONSuccessInfoExpiration int64

const JSONSuccessInfoExpirationNoCache JSONSuccessInfoExpiration = 0

// JSONSuccessOrder <https://vk.com/dev/payments_status>.
type JSONSuccessOrder struct {
	// OrderID, vk description: required идентификатор заказа в системе
	// платежей ВКонтакте.
	OrderID int64 `json:"order_id"`
	// AppOrderID vk description: идентификатор заказа в приложении.
	// Должен быть уникальным для каждого заказа.
	AppOrderID int64 `json:"app_order_id,omitempty"`
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
