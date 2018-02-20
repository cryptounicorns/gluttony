package lua

type Config struct {
	Code         string `validate:"required"`
	FunctionName string `validate:"required"`
	Workers      uint   `validate:"required"`
}
