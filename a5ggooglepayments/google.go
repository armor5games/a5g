package a5ggooglepayments

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"strings"

	"github.com/pkg/errors"
)

// GoogleInappPurchaseReceipt some documentation:
// <https://developer.android.com/google/play/billing/billing_integrate.html#Purchase>,
// <https://developer.android.com/google/play/billing/billing_reference.html>.
// Example:
// {
//    \"orderId\":\"12999763169054705758.1371079406387615\",
//    \"packageName\":\"com.example.app\",
//    \"productId\":"exampleSku\",
//    \"purchaseTime\":1345678900000,
//    \"purchaseState\":0,
//    \"developerPayload\":\"bGoa+V7g/yqDXvKRqq+JTFn4uQZbPiQJo4pf9RzJ\",
//    \"purchaseToken\":\"rojeslcdyyiapnqcynkjyyjh\"
//  }
type GoogleInappPurchaseReceipt struct {
	OrderID          string `json:"orderId"`
	PackageName      string `json:"packageName"`
	ProductID        string `json:"productId"`
	PurchaseTime     uint64 `json:"purchaseTime"`
	PurchaseState    uint64 `json:"purchaseState"`
	DeveloperPayload string `json:"developerPayload,omitempty"`
	PurchaseToken    string `json:"purchaseToken"`
}

// IsValid <https://developer.android.com/google/play/licensing/setting-up.html>.
func IsValid(
	base64EncodedPublicKey, receiptSignature string, saleReceipt []byte) (
	bool, error) {
	publicKey, err := decodePublicKey(base64EncodedPublicKey)
	if err != nil {
		return false, errors.WithStack(err)
	}

	// generate hash value from receipt
	hasher := sha1.New()

	_, err = hasher.Write(saleReceipt)
	if err != nil {
		return false, errors.WithStack(err)
	}

	hashedReceipt := hasher.Sum(nil)

	// decode receiptSignature
	decodedSignature, err := base64.StdEncoding.DecodeString(receiptSignature)
	if err != nil {
		return false, errors.WithStack(err)
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (v *GoogleInappPurchaseReceipt) Validate() error {
	if strings.TrimSpace(v.OrderID) == "" {
		return errors.New("missing google inapp purchase order id")
	}

	if strings.TrimSpace(v.PackageName) == "" {
		return errors.New("missing google inapp purchase ")
	}

	if strings.TrimSpace(v.ProductID) == "" {
		return errors.New("missing google inapp purchase ")
	}

	if v.PurchaseTime < 1 {
		return errors.New("missing google inapp purchase ")
	}

	if strings.TrimSpace(v.PurchaseToken) == "" {
		return errors.New("missing google inapp purchase ")
	}

	return nil
}

func decodePublicKey(base64EncodedPublicKey string) (*rsa.PublicKey, error) {
	decodedPublicKey, err :=
		base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an *rsa.PublicKey")
	}

	return publicKey, nil
}
