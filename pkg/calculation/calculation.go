package calculation

import (
	"strconv"
	"unicode"
)

type Stack[T any] struct {
	lst []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (st *Stack[T]) Push(x T) {
	st.lst = append(st.lst, x)
}

func (st *Stack[T]) Pop() T {
	if len(st.lst) == 0 {
		var zero T
		return zero
	}
	res := st.lst[len(st.lst)-1]
	st.lst = st.lst[:len(st.lst)-1]
	return res
}

func (st *Stack[T]) Back() T {
	if len(st.lst) == 0 {
		var zero T
		return zero
	}
	return st.lst[len(st.lst)-1]
}

func (st *Stack[T]) Len() int {
	return len(st.lst)
}

func toTokens(str string) (res []string, err error) {
	var current = ""
	var prev = ""

	for i, el := range str {
		if unicode.IsSpace(el) {
			continue
		}
		if unicode.IsDigit(el) {
			if current == "-u" {
				res = append(res, current)
				current = ""
			}
			current += string(el)
			prev = ""
		} else {
			if current != "" {
				res = append(res, current)
				current = ""
			}
			if el == '-' && (i == 0 || prev == "(" || prev == "+" || prev == "-" || prev == "*" || prev == "/") {
				current = "-u"
			} else {
				if el == '+' || el == '-' || el == '*' || el == '/' || el == ')' || el == '(' {
					res = append(res, string(el))
				} else {
					return nil, ErrUnknownSymbol
				}
			}
		}
		if current == "" {
			prev = string(el)
		}
	}
	if current != "" {
		res = append(res, current)
	}
	return
}

func priority(op string) int {
	if op == "+" || op == "-" {
		return 1
	} else {
		return 2
	}
}

func toNotation(tokens []string) (res []string, err error) {
	var opers Stack[string]
	for _, el := range tokens {
		if _, err := strconv.Atoi(el); err == nil {
			res = append(res, el)
		} else if el == "+" || el == "-" || el == "*" || el == "/" || el == "-u" {
			for opers.Len() != 0 && opers.Back() != "(" {
				a := priority(opers.Back())
				b := priority(el)
				if a > b || (a == b && el != "-u") {
					if opers.Len() == 0 {
						return nil, ErrInvalidExpression
					} else {
						res = append(res, opers.Pop())
					}
				} else {
					break
				}
			}
			opers.Push(el)
		} else if el == "(" {
			opers.Push(el)
		} else {
			for opers.Len() != 0 && opers.Back() != "(" {
				if opers.Len() == 0 {
					return nil, ErrInvalidExpression
				}
				res = append(res, opers.Pop())
			}
			if opers.Len() == 0 {
				return nil, ErrInvalidExpression
			}
			opers.Pop()
		}
	}
	for opers.Len() != 0 {
		el := opers.Pop()
		if el == "(" {
			return nil, ErrInvalidExpression
		} else {
			res = append(res, el)
		}
	}
	return res, nil
}

func calculateNotation(st []string) (float64, error) {
	var res Stack[float64]
	for _, el := range st {
		if number, err := strconv.ParseFloat(el, 64); err == nil {
			res.Push(number)
		} else {
			if res.Len() == 0 {
				return 0, ErrInvalidExpression
			}
			oper2 := res.Pop()
			if el == "-u" {
				res.Push(-oper2)
				continue
			}

			if res.Len() == 0 {
				return 0, ErrInvalidExpression
			}
			oper1 := res.Pop()
			switch el {
			case "+":
				res.Push(oper1 + oper2)
			case "-":
				res.Push(oper1 - oper2)
			case "*":
				res.Push(oper1 * oper2)
			case "/":
				if oper2 == 0 {
					return 0, ErrDivideByZero
				}
				res.Push(oper1 / oper2)
			}
		}
	}
	if res.Len() != 1 {
		return 0, ErrInvalidExpression
	}
	return res.Pop(), nil
}

func Calc(str string) (float64, error) {
	tokens, err := toTokens(str)
	if err != nil {
		return 0, err
	}
	st, err := toNotation(tokens)
	if err != nil {
		return 0, err
	}
	res, err := calculateNotation(st)
	if err != nil {
		return 0, err
	}
	return res, nil
}
