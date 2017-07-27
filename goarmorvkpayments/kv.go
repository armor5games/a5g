package goarmorvkpayments

import (
	"net/url"

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

	return nil
}

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
