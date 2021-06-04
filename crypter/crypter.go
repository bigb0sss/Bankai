package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	math "math/rand"
)

func RandKeyGen(n int) string {

	// Random Key Generator (128 bit)
	var chars = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	charSet := make([]rune, n)
	for i := range charSet {
		charSet[i] = chars[math.Intn(len(chars))]
	}
	return string(charSet)
}

// Encrpyt: Original Text --> Add IV --> Encrypt with Key --> Base64 Encode
func Encrypt(key []byte, text []byte) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Creating IV
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// AES Encrpytion Process
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	// Base64 Encode
	return base64.URLEncoding.EncodeToString(ciphertext)
}
