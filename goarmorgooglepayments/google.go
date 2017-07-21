package goarmorgooglepayments

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"

	"github.com/pkg/errors"
)

type GoogleInappPurchaseData struct {
	OrderID          string `json:"orderId"`
	PackageName      string `json:"packageName"`
	ProductID        string `json:"productId"`
	PurchaseTime     uint64 `json:"purchaseTime"`
	PurchaseState    uint64 `json:"purchaseState"`
	DeveloperPayload string `json:"developerPayload,omitempty"`
	PurchaseToken    string `json:"purchaseToken"`
}

func DecodePublicKey(file string) (*rsa.PublicKey, error) {
	base64EncodedPublicKey, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(
		string(base64EncodedPublicKey[:]))
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

// {
//    \"orderId\":\"12999763169054705758.1371079406387615\",
//    \"packageName\":\"com.example.app\",
//    \"productId\":"exampleSku\",
//    \"purchaseTime\":1345678900000,
//    \"purchaseState\":0,
//    \"developerPayload\":\"bGoa+V7g/yqDXvKRqq+JTFn4uQZbPiQJo4pf9RzJ\",
//    \"purchaseToken\":\"rojeslcdyyiapnqcynkjyyjh\"
//  }

func IsValid(publicKey *rsa.PublicKey, signature string, receipt []byte) (
	bool, error) {
	// generate hash value from receipt
	hasher := sha1.New()

	if _, err := hasher.Write(receipt); err != nil {
		return false, errors.WithStack(err)
	}

	hashedReceipt := hasher.Sum(nil)

	// decode signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.WithStack(err)
	}

	// verify
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}
