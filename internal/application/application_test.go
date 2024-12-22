package application

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		code       int
	}{
		{
			name:       "code 200",
			expression: `{"expression":"2+2*2"}`,
			code:       200,
		},
		{
			name:       "code 422 unknown symbol",
			expression: `{"expression":"2+a"}`,
			code:       422,
		},
		{
			name:       "code 422 invalid expression",
			expression: `{"expression":"2++4"}`,
			code:       422,
		},
		{
			name:       "code 422 divide by zero",
			expression: `{"expression":"32/0"}`,
			code:       422,
		},
	}

	for _, tt := range tests {
		reader := strings.NewReader(tt.expression)
		req := httptest.NewRequest("POST", "http://localhost:8080/api/v1/calculate", reader)
		w := httptest.NewRecorder()
		CalcHandler(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != tt.code {
			t.Fatalf("case %s return code %d instead of %d", tt.name, resp.StatusCode, tt.code)
		}
	}
}
