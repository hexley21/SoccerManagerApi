package playground_validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func registerDecimal(field reflect.Value) any {
	if valuer, ok := field.Interface().(decimal.Decimal); ok {
		return valuer.String()
	}
	return nil
}

func decimalGTValidator(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}
	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.GreaterThanOrEqual(baseValue)
}
