package app

type ValidationResponse struct {
	Valid  bool
	Issues []string
}

type OperationResult[T any] struct {
	Value T
}
