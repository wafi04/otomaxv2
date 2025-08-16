package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// Hash algorithms
type HashAlgorithm string

const (
	MD5    HashAlgorithm = "md5"
	SHA1   HashAlgorithm = "sha1"
	SHA256 HashAlgorithm = "sha256"
	SHA512 HashAlgorithm = "sha512"
)

// Crypto provides cryptographic operations
type Crypto struct {
	secretKey []byte
}

// NewCrypto creates a new crypto instance
func NewCrypto(secretKey string) *Crypto {
	return &Crypto{
		secretKey: []byte(secretKey),
	}
}

// Password Hashing with bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies password against hash
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// PBKDF2 key derivation
func DeriveKeyPBKDF2(password, salt string, iterations, keyLength int, hashFunc func() hash.Hash) []byte {
	return pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, hashFunc)
}

// Scrypt key derivation
func DeriveKeyScrypt(password, salt string, N, r, p, keyLen int) ([]byte, error) {
	return scrypt.Key([]byte(password), []byte(salt), N, r, p, keyLen)
}

// Hash functions
func Hash(data string, algorithm HashAlgorithm) string {
	switch algorithm {
	case MD5:
		hash := md5.Sum([]byte(data))
		return hex.EncodeToString(hash[:])
	case SHA1:
		hash := sha1.Sum([]byte(data))
		return hex.EncodeToString(hash[:])
	case SHA256:
		hash := sha256.Sum256([]byte(data))
		return hex.EncodeToString(hash[:])
	case SHA512:
		hash := sha512.Sum512([]byte(data))
		return hex.EncodeToString(hash[:])
	default:
		hash := sha256.Sum256([]byte(data))
		return hex.EncodeToString(hash[:])
	}
}

// HMAC signature
func (c *Crypto) GenerateHMAC(data string, algorithm HashAlgorithm) string {
	var mac hash.Hash
	
	switch algorithm {
	case SHA1:
		mac = hmac.New(sha1.New, c.secretKey)
	case SHA256:
		mac = hmac.New(sha256.New, c.secretKey)
	case SHA512:
		mac = hmac.New(sha512.New, c.secretKey)
	default:
		mac = hmac.New(sha256.New, c.secretKey)
	}
	
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMAC verifies HMAC signature
func (c *Crypto) VerifyHMAC(data, signature string, algorithm HashAlgorithm) bool {
	expectedMAC := c.GenerateHMAC(data, algorithm)
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// AES Encryption/Decryption
func (c *Crypto) EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.secretKey[:32]) // Use first 32 bytes for AES-256
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Crypto) DecryptAES(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.secretKey[:32])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}

	nonce, cipherData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Generate random string
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// Generate random bytes
func GenerateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// Generate API Key
func GenerateAPIKey() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	randomStr := GenerateRandomString(16)
	combined := fmt.Sprintf("%s_%s", timestamp, randomStr)
	return base64.URLEncoding.EncodeToString([]byte(combined))
}

// Generate transaction ID
func GenerateTransactionID(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	random := GenerateRandomString(6)
	return fmt.Sprintf("%s_%s_%s", prefix, timestamp, random)
}

// Generate signature for payment gateway
func GeneratePaymentSignature(merchantID, orderID, grossAmount, serverKey string) string {
	signatureString := fmt.Sprintf("%s%s%s%s", merchantID, orderID, grossAmount, serverKey)
	return Hash(signatureString, SHA512)
}

// Verify payment signature
func VerifyPaymentSignature(merchantID, orderID, grossAmount, serverKey, signature string) bool {
	expectedSignature := GeneratePaymentSignature(merchantID, orderID, grossAmount, serverKey)
	return signature == expectedSignature
}

// Base64 encoding/decoding
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// URL-safe Base64 encoding/decoding
func Base64URLEncode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func Base64URLDecode(data string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}