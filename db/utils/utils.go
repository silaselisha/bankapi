package utils

import (
	"math/rand"
	"strings"
)

var letters string = "abcdefghijklmnopqrstuvwxyz"

func RandomAmount(min, max int) int {
	return int(rand.Intn((max - min) + 1))
}

func RandomCurrency() string {
	currency := []string{"USD", "EUR", "KSH", "GBP"}
	return currency[rand.Intn(len(currency))]
}

func RandomString(size int) (string, error) {
	var sb strings.Builder
	for i := 0; i < size; i++ {
		ch := letters[rand.Intn(len(letters))]
		err := sb.WriteByte(ch)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}