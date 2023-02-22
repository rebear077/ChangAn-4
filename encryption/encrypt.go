package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

type Encrypter struct {
	s *symmetry
	a *asymmetry
}

func NewEncrypter() *Encrypter {
	return &Encrypter{
		s: new(symmetry),
		a: new(asymmetry),
	}
}
func (e *Encrypter) Signature(data []byte) []byte {
	signed := rsaSignWithSha256(data)
	return signed
}
func (e *Encrypter) SymEncrypt(data []byte, key []byte) ([]byte, error) {
	encrypted, err := e.s.encode(data, key)
	if err != nil {
		return nil, err
	}
	return encrypted, nil

}
func (e *Encrypter) SymDecrypt(encryptData []byte, key []byte) ([]byte, error) {

	data, err := e.s.decode(encryptData, key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (e *Encrypter) AsymEncrypt(data []byte, key []byte) ([]byte, error) {

	cipherData, err := e.a.encode(data, key)
	if err != nil {
		return nil, err
	}
	return cipherData, nil
}
func (e *Encrypter) AsymDecrypt(encryptData []byte, key []byte) ([]byte, error) {

	cipherData, err := e.a.decode(encryptData, key)
	if err != nil {
		return nil, err
	}
	return cipherData, nil
}

type symmetry struct{}

// 对称加密
func (s *symmetry) encode(data []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	cipherByte := make([]byte, len(data))
	cfb.XORKeyStream(cipherByte, data)
	return cipherByte, nil
}

// 对称解密
func (s *symmetry) decode(cipherByte []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plainByte := make([]byte, len(cipherByte))
	cfbdec.XORKeyStream(plainByte, cipherByte)
	data := plainByte
	return data, nil
}

type asymmetry struct{}

// 非对称加密
func (a *asymmetry) encode(data []byte, publickey []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey := pubInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// 非对称解密
func (a *asymmetry) decode(cipherText []byte, privatekey []byte) ([]byte, error) {
	block, _ := pem.Decode(privatekey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// sha256求哈希值
func rsaSignWithSha256(data []byte) []byte {
	m := sha256.New()
	m.Write(data)
	res := hex.EncodeToString(m.Sum(nil))
	return []byte(res)
}
