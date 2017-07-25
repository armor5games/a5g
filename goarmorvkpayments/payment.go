package goarmorvkpayments

// VKAPIPayment <https://vk.com/dev/payments_callbacks>.
type VKAPIPayment struct {
	// NotificationType (notification_type) тип уведомления.
	// равен "order_status_change" или "order_status_change_test".
	NotificationType string `json:"notificationType"`

	// AppID (app_id) идентификатор приложения.
	AppID int64 `json:"appID"`

	// UserID (user_id) идентификатор пользователя, сделавшего заказ.
	UserID int64 `json:"userID"`

	// ReceiverID (receiver_id) идентификатор получателя заказа
	// (в данный момент совпадает с user_id, но в будущем может отличаться).
	ReceiverID int64 `json:"receiverID"`

	// OrderID (order_id) идентификатор заказа в системе платежей ВКонтакте.
	OrderID int64 `json:"orderID"`

	// Signature (sig) подпись уведомления (см. подробнее в разделе 3. Проверка подписи уведомления).
	Signature int64 `json:"signature"`
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

// VKAPIPaymentOrderStatus <https://vk.com/dev/payments_status>.
type VKAPIPaymentOrderStatus struct {
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
	ItemID int64 `json:"itemID"`

	// ItemTitle (item_title) название товара.
	ItemTitle string `json:"itemTitle"`

	// ItemPhotoURL (item_photo_url) string изображение товара.
	ItemPhotoURL string `json:"itemPhotoURL"`

	// ItemPrice (item_price) стоимость товара.
	ItemPrice string `json:"itemPrice"`
}
