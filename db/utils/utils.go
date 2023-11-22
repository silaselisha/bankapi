package utils

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var letters string = "abcdefghijklmnopqrstuvwxyz"

func RandomAmount(min, max int) int {
	return min + int(rand.Intn((max-min)+1))
}

func RandomCurrency() string {
	currency := []string{"USD", "EUR", "GBP"}
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

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
)

func currencyValidator(curr string) bool {
	switch curr {
	case USD, EUR, GBP:
		{
			return true
		}
	default:
		return false
	}
}

var CurrencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		return currencyValidator(currency)
	}
	return false
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hashing failed %w", err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}