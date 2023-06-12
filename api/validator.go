package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/silaselisha/bank-api/db/utils"
)


var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency, ok := fieldlevel.Field().Interface().(string);  ok {
		return utils.IsSupportedCurrency(currency)
	}

	return false
}