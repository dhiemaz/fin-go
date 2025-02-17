package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"log"
	"time"
)

func GenerateUUID() string {
	uniqueId := uuid.New()
	return uniqueId.String()
}

func GenerateRandomToken() string {
	length := 100
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	token := base64.RawURLEncoding.EncodeToString(randomBytes)
	return token
}

func GenerateCode(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	return prefix + timestamp
}

func GenerateCIF() string {
	// Generate a 12-character random string (you can modify the length or include specific patterns)
	const cifLength = 12
	var cif []byte
	for i := 0; i < cifLength; i++ {
		b := make([]byte, 1)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		// Adding numbers or letters to the CIF
		cif = append(cif, '0'+b[0]%10) // Random digit between 0-9
	}
	return string(cif)
}
