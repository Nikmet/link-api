package request

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)

	if err != nil {
		return payload, err
	}

	fmt.Println(payload)

	return payload, nil
}
