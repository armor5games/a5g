package goarmorpay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
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

func DecodePublickey(file string) (*rsa.PublicKey, error) {
	base64EncodedPublicKey, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(fmt.Errorf("Error during open public key file (%s) : %s\n", file, err.Error()))
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(string(base64EncodedPublicKey[:]))
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key")
	}

	publicKey, _ := publicKeyInterface.(*rsa.PublicKey)

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

func GooglePayVerify(publicKey *rsa.PublicKey, signature string, receipt []byte) (cheater bool, err error) {
	// generate hash value from receipt
	hasher := sha1.New()
	hasher.Write(receipt)
	hashedReceipt := hasher.Sum(nil)

	// decode signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature")
	}

	// verify
	if err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature); err != nil {
		return true, err
	}

	return false, nil
}
