package goarmorvkpayments

import (
	"net/url"
	"strings"

	"github.com/armor5games/goarmor/goarmorapi"
	"github.com/pkg/errors"
)

// VKAPIKV vk payment's key values
type VKAPIKV map[string]string

var ErrVKAPIKVKeyName = errors.New("unknown vk payment's key name")

func NewVKAPIKVByURLValues(keyValues url.Values) VKAPIKV {
	kv := make(VKAPIKV)

	if keyValues == nil {
		return kv
	}

	for k, v := range keyValues {
		if len(v) > 0 {
			kv[k] = v[0]
		}
	}

	return kv
}

func (keyValues VKAPIKV) Validate() error {
	for k := range keyValues {
		switch k {
		default:
			return ErrVKAPIKVKeyName

		case
			"notification_type",
			"app_id",
			"user_id",
			"receiver_id",
			"order_id",
			"sig":
			// VKAPIPayment.

		case "item":
			// VKAPIPaymentInfo and VKAPIPaymentOrder.

		case "lang":
			// VKAPIPaymentInfo.

		case
			"date",
			"item_id",
			"item_photo_url",
			"item_price",
			"item_title",
			"status":
			// VKAPIPaymentOrder.
		}
	}

	if strings.TrimSpace(keyValues["app_id"]) == "" {
		return errors.New(`empty vk payment's "app_id"`)
	}

	if strings.TrimSpace(keyValues["notification_type"]) == "" {
		return errors.New(`empty vk payment's "notification_type"`)
	}

	if strings.TrimSpace(keyValues["order_id"]) == "" {
		return errors.New(`empty vk payment's "order_id"`)
	}

	if strings.TrimSpace(keyValues["receiver_id"]) == "" {
		return errors.New(`empty vk payment's "receiver_id"`)
	}

	if strings.TrimSpace(keyValues["sig"]) == "" {
		return errors.New(`empty vk payment's "sig"`)
	}

	if strings.TrimSpace(keyValues["user_id"]) == "" {
		return errors.New(`empty vk payment's "user_id"`)
	}

	err := VKAPIPaymentNotificationType(keyValues["notification_type"]).Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m VKAPIKV) Item() string             { return m["item"] }
func (m VKAPIKV) NotificationType() string { return m["notification_type"] }
func (m VKAPIKV) OrderID() string          { return m["order_id"] }
func (m VKAPIKV) Sig() string              { return m["sig"] }

func (keyValues VKAPIKV) KV() goarmorapi.KV {
	kv := make(goarmorapi.KV)

	if len(keyValues) == 0 {
		return kv
	}

	for k, v := range keyValues {
		kv[k] = v
	}

	return kv
}

func (keyValues VKAPIKV) KVS() goarmorapi.KVS {
	return goarmorapi.KVS(keyValues)
}
