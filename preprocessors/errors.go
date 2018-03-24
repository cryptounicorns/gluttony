package preprocessors

import (
	"fmt"
)

type ErrUnknownPreprocessorType struct {
	t string
}

func (e *ErrUnknownPreprocessorType) Error() string {
	return fmt.Sprintf(
		"Unknown preprocessor type '%s'",
		e.t,
	)
}
func NewErrUnknownPreprocessorType(t string) error {
	return &ErrUnknownPreprocessorType{t}
}
