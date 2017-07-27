package goarmorvkpayments

import (
	"strings"

	"github.com/pkg/errors"
)

var ErrVKAPIPaymentInfoLanguage = errors.New(`unexpected vk payment's "lang"`)

// VKAPIPaymentInfo <https://vk.com/dev/payments_getitem>.
type VKAPIPaymentInfo struct {
	VKAPIPayment

	// Language (lang) язык пользователя в формате язык_страна.
	// На данный момент поддерживается 4 языка.
	Language string `json:"language"`

	// Item наименование товара, переданное диалоговому окну покупки
	// (см. Параметры диалогового окна платежей)
	Item string `json:"item"`
}

// VKAPIPaymentInfoLanguage allowable languages for VKAPIPaymentInfo.Language
type VKAPIPaymentInfoLanguage string

const (
	VKAPIPaymentInfoLanguageRURU VKAPIPaymentInfoLanguage = "ru_RU"
	VKAPIPaymentInfoLanguageUKUA VKAPIPaymentInfoLanguage = "uk_UA"
	VKAPIPaymentInfoLanguageBEBY VKAPIPaymentInfoLanguage = "be_BY"
	VKAPIPaymentInfoLanguageENUS VKAPIPaymentInfoLanguage = "en_US"
)

func (t VKAPIPaymentInfoLanguage) String() string {
	return string(t)
}

func (p *VKAPIPaymentInfo) Validate() error {
	err := p.VKAPIPayment.Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	switch VKAPIPaymentNotificationType(p.NotificationType) {
	default:
		return ErrVKAPIPaymentUnexpectedNotificationType
	case
		VKAPIPaymentNotificationTypeGetItem,
		VKAPIPaymentNotificationTypeGetItemTest:
	}

	switch VKAPIPaymentInfoLanguage(p.Language) {
	default:
		return ErrVKAPIPaymentInfoLanguage
	case
		VKAPIPaymentInfoLanguageBEBY,
		VKAPIPaymentInfoLanguageENUS,
		VKAPIPaymentInfoLanguageRURU,
		VKAPIPaymentInfoLanguageUKUA:
	}

	if strings.TrimSpace(p.Item) == "" {
		return ErrVKAPIPaymentItemEmpty
	}

	return nil
}

func (kv VKAPIKV) VKAPIPaymentInfo() (*VKAPIPaymentInfo, error) {
	p, err := kv.kvAPIPayment()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	paymentInfo := &VKAPIPaymentInfo{VKAPIPayment: p}

	for k, v := range kv {
		switch k {
		case "lang":
			paymentInfo.Language = v

		case "item":
			paymentInfo.Item = v
		}
	}

	return paymentInfo, nil
}
