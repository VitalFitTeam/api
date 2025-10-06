package otp

import (
	"fmt"
	"math/rand"
	"time"
)

const defaultAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(length int) (string, error) {

	rand.NewSource(time.Now().UnixNano())

	code := make([]rune, length)
	alphabetRunes := []rune(defaultAlphabet)
	alphabetSize := len(alphabetRunes)

	if alphabetSize == 0 {
		return "", fmt.Errorf("el alfabeto no puede estar vacío")
	}

	for i := 0; i < length; i++ {
		// Selecciona un carácter aleatorio del alfabeto
		randomIndex := rand.Intn(alphabetSize)
		code[i] = alphabetRunes[randomIndex]
	}

	return string(code), nil
}
