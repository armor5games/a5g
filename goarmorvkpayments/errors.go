package goarmorvkpayments

import "errors"

// <https://vk.com/dev/payments_errors>
var (
	ErrDBTemporaryError = JSONError{
		Code:     2,
		Critical: false,
		Err:      errors.New("temporary database error")}

	ErrSignatureMismatch = JSONError{
		Code:     10,
		Critical: true,
		Err:      errors.New("mismatching of the calculated and transmitted signature")}

	ErrQueryParameters = JSONError{
		Code:     11,
		Critical: true,
		Err:      errors.New("query parameters do not meet the specification")}

	ErrProductNotExist = JSONError{
		Code:     20,
		Critical: true,
		Err:      errors.New("product does not exist")}

	ErrProductOutOfStock = JSONError{
		Code:     21,
		Critical: true,
		Err:      errors.New("product is out of stock")}

	ErrUserNotExist = JSONError{
		Code:     22,
		Critical: true,
		Err:      errors.New("user does not exist")}

	ErrUnknownInternalServerError = JSONError{
		Code:     100,
		Critical: true,
		Err:      errors.New("unknown internal server error")}

	ErrUnsupportedPaymentNotificationType = JSONError{
		Code:     101,
		Critical: true,
		Err:      errors.New("unsupported payment notification type")}
)
