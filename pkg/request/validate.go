package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateBody[T any](p T) error {
	validate := validator.New()
	err := validate.Struct(p)

	if err != nil {
		return err
	}
	fmt.Println("Valid")
	return nil
}
