package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

// pkcs7Pad 测试用填充（生产解密路径不需要，仅用于构造测试密文）
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func encryptForTest(t *testing.T, key, iv, plaintext []byte) string {
	t.Helper()
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	padded := pkcs7Pad(plaintext, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(padded))
	mode.CryptBlocks(out, padded)
	return base64.StdEncoding.EncodeToString(out)
}

func TestDecryptPhoneNumber(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(iv); err != nil {
		t.Fatal(err)
	}

	plaintext := `{"phoneNumber":"13800001111","purePhoneNumber":"13800001111","countryCode":"86"}`
	encryptedData := encryptForTest(t, key, iv, []byte(plaintext))
	sessionKey := base64.StdEncoding.EncodeToString(key)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	phone, err := DecryptPhoneNumber(sessionKey, encryptedData, ivStr)
	if err != nil {
		t.Fatalf("DecryptPhoneNumber failed: %v", err)
	}
	if phone != "13800001111" {
		t.Fatalf("phone = %q, want 13800001111", phone)
	}
}

func TestDecryptPhoneNumberInvalidSessionKey(t *testing.T) {
	if _, err := DecryptPhoneNumber("not-base64!!", "AAAA", "AAAA"); err == nil {
		t.Fatal("expected error for invalid session_key")
	}
}
