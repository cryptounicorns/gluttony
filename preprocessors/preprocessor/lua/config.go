package lua

type Config struct {
	Code         string `validator:"required"`
	FunctionName string `validator:"required"`
	Workers      uint   `validator:"required"`
}
