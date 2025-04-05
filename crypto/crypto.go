package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/dexlabsio/garlic/errors"
)

type Manager interface {
	Encrypt([]byte) (string, error)
	Decrypt(string) ([]byte, error)
}

type Crypto struct {
	config *Config
}

func New(cfg *Config) *Crypto {
	return &Crypto{cfg}
}

// Encrypt takes a random text and uses AES to encrypt the text.
func (c *Crypto) Encrypt(content []byte) (string, error) {
	block, err := aes.NewCipher(c.generateAESKey())
	if err != nil {
		return "", errors.Propagate(err, "failed to generate crypto cypher")
	}

	// Generate a random initialization vector.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Propagate(err, "failed to generate random initialization vector")
	}

	// Pad the content to a multiple of the block size
	padding := aes.BlockSize - (len(content) % aes.BlockSize)
	paddedContent := append(content, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// Create a new AES cipher block mode
	ciphertext := make([]byte, aes.BlockSize+len(paddedContent))
	ivStart := ciphertext[:aes.BlockSize]
	copy(ivStart, iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], paddedContent)

	// Encode the encrypted text as base64 for storage or transmission
	encodedText := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedText, nil
}

func (c *Crypto) Decrypt(encodedText string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return nil, errors.Propagate(err, "failed to decode base64 string")
	}

	block, err := aes.NewCipher(c.generateAESKey())
	if err != nil {
		return nil, errors.Propagate(err, "failed to create new cipher")
	}

	if len(encrypted) < aes.BlockSize {
		return nil, errors.New(errors.KindSystemError, "invalid ciphertext length")
	}

	iv := encrypted[:aes.BlockSize]
	ciphertext := encrypted[aes.BlockSize:]

	// Create a new byte slice for decrypted data
	decrypted := make([]byte, len(ciphertext))

	// Perform the decryption
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	// Unpad the decrypted password
	padding := int(decrypted[len(decrypted)-1])
	decrypted = decrypted[:len(decrypted)-padding]

	return decrypted, nil
}

// generateAESKey creates a 32-byte compatible AES Key from a random
// configured MasterKey string.
func (c *Crypto) generateAESKey() []byte {
	// Convert the random string to bytes
	randomBytes := []byte(c.config.MasterKey)

	// Generate a 32-byte AES key using SHA-256 hash function
	hash := sha256.Sum256(randomBytes)
	aesKey := hash[:]

	return aesKey
}
