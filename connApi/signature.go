package receive

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// // 签名
// func sha256withRSASignature(data []byte, keyBytes []byte) (string, error) {
// 	h := sha256.New()
// 	h.Write(data)
// 	hashed := h.Sum(nil)
// 	block, _ := pem.Decode(keyBytes)
// 	if block == nil {
// 		return "", errors.New("private key error")
// 	}
// 	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		fmt.Println("ParsePKCS8PrivateKey err", err)
// 		return "", err
// 	}

// 	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
// 	if err != nil {
// 		fmt.Printf("Error from signing: %s\n", err)
// 		return "", err
// 	}

// 	return base64.StdEncoding.EncodeToString(signature), nil
// }

// 验证
func rsaVerySignWithSha256(data []byte, signData string, keyBytes []byte) (bool, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return false, errors.New("public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(data)
	signature, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, err
	}
	return true, nil
}
