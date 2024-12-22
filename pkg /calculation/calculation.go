package calculation

import (
	"errors"
	"strconv"
	"strings"
)

// Calc принимает строковое выражение и возвращает результат или ошибку
func Calc(expression string) (float64, error) {
	// Убираем пробелы для корректной обработки
	expression = strings.ReplaceAll(expression, " ", "")

	// Используем стек для чисел и операций
	numStack := []float64{}
	opStack := []rune{}

	// Вспомогательная функция для выполнения операций
	applyOp := func() error {
		if len(numStack) < 2 || len(opStack) == 0 {
			return errors.New("ошибка в выражении")
		}
		b := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]
		a := numStack[len(numStack)-1]
		numStack = numStack[:len(numStack)-1]

		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		switch op {
		case '+':
			numStack = append(numStack, a+b)
		case '-':
			numStack = append(numStack, a-b)
		case '*':
			numStack = append(numStack, a*b)
		case '/':
			if b == 0 {
				return errors.New("деление на ноль")
			}
			numStack = append(numStack, a/b)
		default:
			return errors.New("неизвестная операция")
		}
		return nil
	}

	// Приоритеты операций
	precedence := func(op rune) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		default:
			return 0
		}
	}

	// Основной цикл для парсинга и вычисления выражения
	for i := 0; i < len(expression); i++ {
		ch := rune(expression[i])

		if ch >= '0' && ch <= '9' || ch == '.' {
			j := i
			for j < len(expression) && (expression[j] >= '0' && expression[j] <= '9' || expression[j] == '.') {
				j++
			}
			num, err := strconv.ParseFloat(expression[i:j], 64)
			if err != nil {
				return 0, errors.New("ошибка преобразования числа")
			}
			numStack = append(numStack, num)
			i = j - 1
		} else if ch == '(' {
			opStack = append(opStack, ch)
		} else if ch == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			if len(opStack) == 0 {
				return 0, errors.New("несбалансированные скобки")
			}
			opStack = opStack[:len(opStack)-1]
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(ch) {
				if err := applyOp(); err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, ch)
		} else {
			return 0, errors.New("недопустимый символ в выражении")
		}
	}

	// Выполняем оставшиеся операции
	for len(opStack) > 0 {
		if err := applyOp(); err != nil {
			return 0, err
		}
	}

	if len(numStack) != 1 {
		return 0, errors.New("ошибка в выражении")
	}

	return numStack[0], nil
}


