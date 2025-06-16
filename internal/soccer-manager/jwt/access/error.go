package access

import "fmt"

func NewInvalidTokenError(userData Data) error {
	return fmt.Errorf("invalid access token: %v", userData)
}