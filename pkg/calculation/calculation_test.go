package calculation_test

import (
	"testing"

	"github.com/kleo-53/web_calc_go/pkg/calculation"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "based",
			expression:     "1+1",
			expectedResult: 2.0,
		}, {
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8.0,
		}, {
			name:           "nonpriority",
			expression:     "2+2*2",
			expectedResult: 6.0,
		}, {
			name:           "divide",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}
	for _, tc := range testCasesSuccess {
		t.Run(tc.name, func(t *testing.T) {
			val, err := calculation.Calc(tc.expression)
			if err != nil {
				t.Fatalf("success case %s returns error", tc.expression)
			}
			if val != tc.expectedResult {
				t.Fatalf("%f should be equal %f", val, tc.expectedResult)
			}
		})
	}
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "div by zero",
			expression:  "123/0",
			expectedErr: calculation.ErrDivideByZero,
		}, {
			name:        "letters",
			expression:  "qwerty",
			expectedErr: calculation.ErrUnknownSymbol,
		}, {
			name:        "brackets",
			expression:  "2*1+2)",
			expectedErr: calculation.ErrInvalidExpression,
		}, {
			name:        "unfinished expression",
			expression:  "2+",
			expectedErr: calculation.ErrInvalidExpression,
		},
		{
			name:        "empty",
			expression:  "",
			expectedErr: calculation.ErrInvalidExpression,
		},
	}
	for _, tc := range testCasesFail {
		t.Run(tc.name, func(t *testing.T) {
			val, err := calculation.Calc(tc.expression)
			if err == nil {
				t.Fatalf("fail case %s returns %f", tc.expression, val)
			}
			if err != tc.expectedErr {
				t.Fatalf("%s should be equal %s", err.Error(), tc.expectedErr.Error())
			}
		})
	}
}
