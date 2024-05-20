package util

import (
	"crypto/aes"
	"crypto/cipher"
	crypt "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	StatusCode int
	Message    string
	Success    bool
}

// Define the alphanumeric character set
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateUniqueAlphaNumeric(length int) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the generated string
	result := make([]byte, length)

	// Fill the byte slice with random characters from the charset
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return
	return string(result)
}

// HashPassword generates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a password with its hashed version.
// It returns true if the password matches the hash, otherwise false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Handle Errors

func HandleErrors(errorType error) *ErrorResponse {
	errorMessage := errorType.Error()
	if strings.Contains(errorMessage, "validation") {
		errResp := ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    strings.TrimPrefix(errorMessage, "validation: "),
			Success:    false,
		}
		return &errResp
	} else if strings.Contains(errorMessage, "not found") {
		errResp := ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    errorMessage,
			Success:    false,
		}
		return &errResp
	} else {
		errResp := ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    errorMessage,
			Success:    false,
		}
		return &errResp
	}
}

func PowerOf10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func AESEncrypt(data, key string) (string, error) {
	// Convert the key to 32 bytes (AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("key must be 32 bytes long for AES-256 encryption")
	}

	// Create a new AES block cipher using the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher block: %v", err)
	}

	// Generate a new IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(crypt.Reader, iv); err != nil {
		return "", fmt.Errorf("error generating IV: %v", err)
	}

	// Encrypt the data
	encryptedData := make([]byte, len(data))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(encryptedData, []byte(data))

	// Prepend the IV to the encrypted data
	encryptedDataWithIV := append(iv, encryptedData...)

	// Encode the encrypted data with IV to base64
	encodedData := base64.StdEncoding.EncodeToString(encryptedDataWithIV)

	return encodedData, nil
}

func AESDecrypt(encodedData, key string) (string, error) {
	// Convert the key to 32 bytes (AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("key must be 32 bytes long for AES-256 encryption")
	}

	// Decode the base64 encoded data
	encryptedDataWithIV, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 encoded data: %v", err)
	}

	// Extract the IV from the encrypted data
	if len(encryptedDataWithIV) < aes.BlockSize {
		return "", fmt.Errorf("encrypted data too short")
	}
	iv := encryptedDataWithIV[:aes.BlockSize]
	encryptedData := encryptedDataWithIV[aes.BlockSize:]

	// Create a new AES block cipher using the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher block: %v", err)
	}

	// Decrypt the data
	decryptedData := make([]byte, len(encryptedData))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(decryptedData, encryptedData)

	return string(decryptedData), nil
}
