package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// phoneInfo 微信手机号解密后的 JSON 结构（仅取所需字段）
type phoneInfo struct {
	PhoneNumber string `json:"phoneNumber"`
}

// DecryptPhoneNumber 用 session_key 解密微信 getPhoneNumber 返回的加密数据
// AES-128-CBC，key=session_key（base64 解码），iv=iv（base64 解码），PKCS7 填充
func DecryptPhoneNumber(sessionKey, encryptedData, iv string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return "", errors.New("session_key 解码失败")
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", errors.New("iv 解码失败")
	}
	dataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", errors.New("encrypted_data 解码失败")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	if len(dataBytes) == 0 || len(dataBytes)%block.BlockSize() != 0 {
		return "", errors.New("加密数据长度无效")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	plain := make([]byte, len(dataBytes))
	mode.CryptBlocks(plain, dataBytes)

	plain, err = pkcs7Unpad(plain, block.BlockSize())
	if err != nil {
		return "", err
	}

	var info phoneInfo
	if err := json.Unmarshal(plain, &info); err != nil {
		return "", errors.New("解密结果解析失败")
	}
	if info.PhoneNumber == "" {
		return "", errors.New("未获取到手机号")
	}
	return info.PhoneNumber, nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据为空")
	}
	padLen := int(data[length-1])
	if padLen == 0 || padLen > blockSize || padLen > length {
		return nil, errors.New("填充数据无效")
	}
	if !bytes.Equal(data[length-padLen:], bytes.Repeat([]byte{byte(padLen)}, padLen)) {
		return nil, errors.New("填充数据校验失败")
	}
	return data[:length-padLen], nil
}
