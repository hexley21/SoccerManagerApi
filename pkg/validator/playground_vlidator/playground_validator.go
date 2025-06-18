package playground_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/pkg/logger"
	"github.com/shopspring/decimal"
)

type playgroundValidator struct {
	validator *validator.Validate
}

func New(logger logger.Logger) *playgroundValidator {
	validate := validator.New()

	if err := validate.RegisterValidation("password", passwordValidator); err != nil {
		logger.Fatalf("failed to register password validator: %v", err)
	}

	if err := validate.RegisterValidation("username", usernameValidator); err != nil {
		logger.Fatalf("failed to register username validator: %v", err)
	}

	if err := validate.RegisterValidation("userrole", userRoleValidator); err != nil {
		logger.Fatalf("failed to register userrole validator: %v", err)
	}

	if err := validate.RegisterValidation("playerpos", playerPositionCodeValidator); err != nil {
		logger.Fatalf("failed to register playerpos validator: %v", err)
	}

	if err := validate.RegisterValidation("localecode", localeCodeValidator); err != nil {
		logger.Fatalf("failed to register localecode validator: %v", err)
	}

	if err := validate.RegisterValidation("countrycode", countryCodeValidator); err != nil {
		logger.Fatalf("failed to register countrycode validator: %v", err)
	}

	validate.RegisterCustomTypeFunc(registerDecimal, decimal.Decimal{})
	if err := validate.RegisterValidation("dgte", decimalGTValidator); err != nil {
		logger.Fatal("failed to register dgte validator")
	}

	return &playgroundValidator{validator: validate}
}

func (v *playgroundValidator) Validate(i any) error {
	err := v.validator.Struct(i)
	if err != nil {
		return err
	}

	return nil
}

func (v *playgroundValidator) ValidateVar(i any, tag string) error {
	err := v.validator.Var(i, tag)
	if err != nil {
		return err
	}

	return nil
}

func passwordValidator(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[\x21-\x7E]{8,36}$`).MatchString(fl.Field().String())
}

func usernameValidator(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,15}$`).MatchString(fl.Field().String())
}

func userRoleValidator(fl validator.FieldLevel) bool {
	return domain.UserRole(fl.Field().String()).Valid()
}

func playerPositionCodeValidator(fl validator.FieldLevel) bool {
	return domain.PlayerPositionCode(fl.Field().String()).Valid()
}

func localeCodeValidator(fl validator.FieldLevel) bool {
	return domain.LocaleCode(fl.Field().String()).Valid()
}

func countryCodeValidator(fl validator.FieldLevel) bool {
	return domain.CountryCode(fl.Field().String()).Valid()
}
