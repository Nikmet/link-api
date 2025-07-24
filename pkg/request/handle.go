package request

import (
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		return nil, err
	}
	err = ValidateBody(body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
