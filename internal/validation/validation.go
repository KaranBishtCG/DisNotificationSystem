package validation

import "fmt"

type ValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s : %s", e.Field, e.Msg)
}
