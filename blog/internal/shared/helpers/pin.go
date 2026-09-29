package helpers

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "strconv"

    "golang.org/x/crypto/bcrypt"
)

// GeneratePIN generates a random numeric PIN of specified length
func GeneratePIN(length int) (string, error) {
    if length <= 0 {
        return "", fmt.Errorf("length must be greater than 0")
    }
    
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    
    pin := ""
    for _, b := range bytes {
        pin += strconv.Itoa(int(b % 10))
        if len(pin) >= length {
            break
        }
    }
    
    return pin, nil
}

// HashPINSecure hashes a PIN using bcrypt (slower, more secure)
func HashPINSecure(pin string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

// VerifyPINSecure verifies a PIN against bcrypt hash
func VerifyPINSecure(pin, hashedPIN string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPIN), []byte(pin))
}

// HashPINFast hashes a PIN using SHA256 (faster, less secure)
// Use this for short-lived OTPs (5-10 minutes)
func HashPINFast(pin string) string {
    hash := sha256.Sum256([]byte(pin))
    return base64.StdEncoding.EncodeToString(hash[:])
}

// VerifyPINFast verifies a PIN against SHA256 hash
func VerifyPINFast(pin, hashedPIN string) bool {
    return HashPINFast(pin) == hashedPIN
}