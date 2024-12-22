package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("incorrect expression")
	ErrDivideByZero      = errors.New("division by zero")
	ErrUnknownSymbol     = errors.New("unknown symbol")
)
