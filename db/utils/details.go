package utils

import (
	"math/rand"
	"strings"
	"time"
)

var letters string = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func randomName(n int) string {
	var s strings.Builder
    l := len(letters)

	for i := 0; i < n; i++ {
		c := letters[rand.Intn(l)]
		s.WriteByte(c)
	}

	return s.String()
}

func GenerateFirstName() string {
	return randomName(6)
}

func GenerateLastName() string {
	return randomName(6)
}

func GenerateAmount() int64 {
	return RandomInteger(0, 1000000)
}

func GenerateCurrency() string {
	currency := []string{"USD", "EUR", "CAD"}
	return currency[rand.Intn(len(currency))]
}

func GenerateGender() string {
	gender := []string{"male", "female"}
	return gender[rand.Intn(len(gender))]
}