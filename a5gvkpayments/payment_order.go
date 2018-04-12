package a5gvkpayments

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrVKAPIPaymentOrderStatusUnexpected = errors.New(`unexpected vk payment's order "status"`)

// VKAPIPaymentOrder <https://vk.com/dev/payments_status>.
type VKAPIPaymentOrder struct {
	VKAPIPayment

	// Date дата создания заказа (в формате Unixtime).
	Date int64 `json:"date"`

	// Status новый статус заказа.
	// Возможные значения:
	// chargeable — заказ готов к оплате. Необходимо оформить заказ
	// пользователю внутри приложения. В случае ответа об успехе
	// платёжная система зачислит голоса на счёт приложения. Если в
	// ответ будет получено сообщение об ошибке, заказ отменяется.
	Status string `json:"status"`

	// Item наименование товара, переданное диалоговому окну покупки
	// (см. Параметры диалогового окна платежей)
	Item string `json:"item"`

	// ItemID (item_id) идентификатор товара в приложении.
	ItemID string `json:"itemId"`

	// ItemTitle (item_title) название товара.
	ItemTitle string `json:"itemTitle"`

	// ItemPhotoURL (item_photo_url) string изображение товара.
	ItemPhotoURL string `json:"itemPhotoUrl"`

	// ItemPrice (item_price) стоимость товара.
	ItemPrice string `json:"itemPrice"`
}

// VKAPIPaymentOrderStatus allowable values for VKAPIPaymentOrder.Status
type VKAPIPaymentOrderStatus string

const (
	VKAPIPaymentOrderStatusChargeable VKAPIPaymentOrderStatus = "chargeable"
)

func (s VKAPIPaymentOrderStatus) String() string {
	return string(s)
}

func (p *VKAPIPaymentOrder) Validate() error {
	err := p.VKAPIPayment.Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	switch VKAPIPaymentNotificationType(p.NotificationType) {
	default:
		return ErrVKAPIPaymentUnexpectedNotificationType
	case
		VKAPIPaymentNotificationTypeOrderStatusChange,
		VKAPIPaymentNotificationTypeOrderStatusChangeTest:
	}

	switch VKAPIPaymentOrderStatus(p.Status) {
	default:
		return ErrVKAPIPaymentOrderStatusUnexpected
	case VKAPIPaymentOrderStatusChargeable:
	}

	if p.Date < 1 {
		return errors.New(`unexpected vk payment's "date"`)
	}

	if strings.TrimSpace(p.Status) == "" {
		return errors.New(`empty vk payment's "status"`)
	}

	if strings.TrimSpace(p.Item) == "" {
		return ErrVKAPIPaymentItemEmpty
	}

	if strings.TrimSpace(p.ItemID) == "" {
		return errors.New(`empty vk payment's "item_id"`)
	}

	if strings.TrimSpace(p.ItemTitle) == "" {
		return errors.New(`empty vk payment's "item_title"`)
	}

	if strings.TrimSpace(p.ItemPhotoURL) == "" {
		return errors.New(`empty vk payment's "item_photo_url"`)
	}

	if strings.TrimSpace(p.ItemPrice) == "" {
		return errors.New(`empty vk payment's "item_price"`)
	}

	return nil
}

func (kv VKAPIKV) VKAPIPaymentOrder() (*VKAPIPaymentOrder, error) {
	p, err := kv.kvAPIPayment()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	paymentOrderStatus := &VKAPIPaymentOrder{VKAPIPayment: p}

	var i int64

	for k, v := range kv {
		switch k {
		case "date":
			i, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			paymentOrderStatus.Date = i

		case "status":
			paymentOrderStatus.Status = v

		case "item":
			paymentOrderStatus.Item = v

		case "item_id":
			paymentOrderStatus.ItemID = v

		case "item_title":
			paymentOrderStatus.ItemTitle = v

		case "item_photo_url":
			paymentOrderStatus.ItemPhotoURL = v

		case "item_price":
			paymentOrderStatus.ItemPrice = v
		}
	}

	return paymentOrderStatus, nil
}
