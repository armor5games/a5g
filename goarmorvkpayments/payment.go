package goarmorvkpayments

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrVKAPIPaymentUnexpectedNotificationType = errors.New(`unexpected vk payment's "notification_type"`)
	ErrVKAPIPaymentItemEmpty                  = errors.New(`empty vk payment's "item"`)
)

// VKAPIPayment <https://vk.com/dev/payments_callbacks>.
type VKAPIPayment struct {
	// NotificationType (notification_type) тип уведомления.
	// равен "order_status_change" или "order_status_change_test".
	NotificationType string `json:"notificationType"`

	// AppID (app_id) идентификатор приложения.
	AppID int64 `json:"appId"`

	// UserID (user_id) идентификатор пользователя, сделавшего заказ.
	UserID int64 `json:"userId"`

	// ReceiverID (receiver_id) идентификатор получателя заказа
	// (в данный момент совпадает с user_id, но в будущем может отличаться).
	ReceiverID int64 `json:"receiverId"`

	// OrderID (order_id) идентификатор заказа в системе платежей ВКонтакте.
	OrderID int64 `json:"orderId"`

	// Signature (sig) подпись уведомления (см. подробнее в разделе 3. Проверка подписи уведомления).
	Signature string `json:"signature"`
}

// VKAPIPaymentNotificationType allowable values for VKAPIPayment.NotificationType
type VKAPIPaymentNotificationType string

const (
	// VKAPIPaymentNotificationTypeGetItem получение информации о товаре.
	VKAPIPaymentNotificationTypeGetItem     VKAPIPaymentNotificationType = "get_item"
	VKAPIPaymentNotificationTypeGetItemTest VKAPIPaymentNotificationType = "get_item_test"

	// VKAPIPaymentNotificationTypeOrderStatusChange изменение статуса заказа.
	VKAPIPaymentNotificationTypeOrderStatusChange     VKAPIPaymentNotificationType = "order_status_change"
	VKAPIPaymentNotificationTypeOrderStatusChangeTest VKAPIPaymentNotificationType = "order_status_change_test"

	// VKAPIPaymentNotificationTypeGetSubscription получение информации о подписке.
	VKAPIPaymentNotificationTypeGetSubscription VKAPIPaymentNotificationType = "get_subscription"

	// VKAPIPaymentNotificationTypeSubscriptionStatusChange изменение статуса подписки.
	VKAPIPaymentNotificationTypeSubscriptionStatusChange VKAPIPaymentNotificationType = "subscription_status_change"
)

func (t VKAPIPaymentNotificationType) String() string {
	return string(t)
}

func (p *VKAPIPayment) Validate() error {
	switch VKAPIPaymentNotificationType(p.NotificationType) {
	default:
		return ErrVKAPIPaymentUnexpectedNotificationType
	case
		VKAPIPaymentNotificationTypeGetItem,
		VKAPIPaymentNotificationTypeGetItemTest,
		VKAPIPaymentNotificationTypeOrderStatusChange,
		VKAPIPaymentNotificationTypeOrderStatusChangeTest,
		VKAPIPaymentNotificationTypeGetSubscription,
		VKAPIPaymentNotificationTypeSubscriptionStatusChange:
	}

	if p.AppID < 0 {
		return errors.New(`unexpected vk payment's "app_id"`)
	}

	if p.UserID < 0 {
		return errors.New(`unexpected vk payment's "user_id"`)
	}

	if p.ReceiverID == 0 {
		return errors.New(`unexpected vk payment's "receiver_id"`)
	}

	if p.OrderID < 0 {
		return errors.New(`unexpected vk payment's "order_id"`)
	}

	if strings.TrimSpace(p.Signature) == "" {
		return errors.New("empty vk payment's signature (sig)")
	}

	return nil
}

func (kv VKAPIKV) kvAPIPayment() (VKAPIPayment, error) {
	var p VKAPIPayment

	var (
		i   int64
		err error
	)

	for k, v := range kv {
		switch k {
		case "notification_type":
			p.NotificationType = v

		case "app_id":
			i, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return p, errors.WithStack(err)
			}
			p.AppID = i

		case "user_id":
			i, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return p, errors.WithStack(err)
			}
			p.UserID = i

		case "receiver_id":
			i, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return p, errors.WithStack(err)
			}
			p.ReceiverID = i

		case "order_id":
			i, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return p, errors.WithStack(err)
			}
			p.OrderID = i

		case "sig":
			p.Signature = v
		}
	}

	return p, nil
}

// func (p *VKAPIPayment) setFieldInt64(fieldName, fieldValue string) error {
// 	i, err := strconv.ParseInt(fieldValue, 10, 64)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	switch fieldName {
// 	default:
// 		return ErrVKAPIPaymentUnknownFieldName
// 	case "app_id":
// 		p.AppID = i
// 	case "user_id":
// 		p.UserID = i
// 	case "receiver_id":
// 		p.ReceiverID = i
// 	case "order_id":
// 		p.OrderID = i
// 	}
// 	return nil
// }
