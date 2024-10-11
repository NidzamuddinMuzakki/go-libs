package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func GenerateSalt(length int) (string, error) {
	rand.NewSource(time.Now().UnixNano())
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = byte(rand.Intn(128))
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password string, salt string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password + salt))
	if err != nil {
		return "", err
	}
	hashedPassword := hex.EncodeToString(hash.Sum(nil))
	return salt + ":" + hashedPassword, nil
}
