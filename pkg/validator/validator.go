package validator

type Validator interface {
	Validate(i any) error
	ValidateVar(i any, tag string) error
}
