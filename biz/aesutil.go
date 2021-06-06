package biz

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// aes加密
func AesEncrypt(origData, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	originBytes := []byte(origData)
	originBytes = pKCS5Padding(originBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	crypted := make([]byte, len(originBytes))
	blockMode.CryptBlocks(crypted, originBytes)
	return base64.URLEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(crypted, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	cryptedBytes, err := base64.URLEncoding.DecodeString(crypted)
	origData := make([]byte, len(cryptedBytes))
	if nil != err {
		return "", err
	}
	blockMode.CryptBlocks(origData, cryptedBytes)
	origData = pKCS5UnPadding(origData)
	return string(origData), nil
}