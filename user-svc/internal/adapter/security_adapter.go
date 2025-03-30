package adapter

import (
	"be-yourmoments/user-svc/internal/helper/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

type SecurityAdapter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type securityAdapter struct {
	keySecret []byte
}

func NewSecurityAdapter() SecurityAdapter {
	keySecretStr := utils.GetEnv("KEY_SECRET_ENCRYPT")
	keySecret := []byte(keySecretStr)

	return &securityAdapter{
		keySecret: keySecret,
	}
}

func (c *securityAdapter) Encrypt(plaintext string) (string, error) {
	key := c.keySecret
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

func (c *securityAdapter) Decrypt(ciphertext string) (string, error) {
	key := c.keySecret
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fiber.ErrBadRequest
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fiber.ErrBadRequest
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fiber.ErrBadRequest
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// TO DO SANITATOR
func (c *securityAdapter) SanitiseStruct(input interface{}) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		log.Println("Input must be a pointer to a struct")
		return
	}
	val = val.Elem()

	policy := bluemonday.UGCPolicy()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			sanitised := policy.Sanitize(field.String())
			field.SetString(sanitised)

		case reflect.Struct:
			// Rekursif jika nested struct
			c.SanitiseStruct(field.Addr().Interface())
		}
	}
}
