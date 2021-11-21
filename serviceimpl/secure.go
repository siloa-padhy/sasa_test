package serviceimpl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func encrypt(maintext string, mainkey string) (string, error) {
	fmt.Println("ernter encypt")
	key := []byte(mainkey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "nil", err
	}
	text := []byte(maintext)
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "nil", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return string(ciphertext), nil
}

func decrypt(text []byte, mainkey string) (string, error) {
	key := []byte(mainkey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "nil", err
	}
	if len(text) < aes.BlockSize {
		return "nil", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		fmt.Println("decrtypting err")
		return "", err
	}
	return string(data), nil
}