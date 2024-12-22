package calculation

import (
	"testing"
)

func TestCalc(t *testing.T) {
	t.Run("Success cases", func(t *testing.T) {
		testCases := []struct {
			name     string
			expr     string
			expected float64
		}{
			{"Simple addition", "2+2", 4},
			{"Complex expression", "(3+2)*4-5/(1+1)", 19.5},
			{"Division", "10/2", 5},
			{"Subtraction", "7-3", 4},
			{"Multiplication", "3*3", 9},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := Calc(tc.expr)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("Error cases", func(t *testing.T) {
		testCases := []struct {
			name     string
			expr     string
			expected error
		}{
			{"Empty expression", "", ErrEmptyExpression},
			{"Short expression", "2", ErrShortExpression},
			{"Division by zero", "10/0", ErrDevisionByZero},
			{"No opening parenthesis", "3+2)*4", ErrNoOpeningParenthesis},
			{"No closing parenthesis", "(3+2*4", ErrNoClosingParenthesis},
			{"Invalid character", "3+2*4$", ErrInvalidExpression},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := Calc(tc.expr)
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err != tc.expected {
					t.Errorf("expected error %v, got %v", tc.expected, err)
				}
			})
		}
	})
}
