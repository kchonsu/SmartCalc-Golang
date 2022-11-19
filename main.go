package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	str = strings.ReplaceAll(str, " ", "")
	str = strings.TrimSpace(str)
	input := []rune(str)
	err := CheckInput(input)
	if err != nil {
		fmt.Println(err)
	} else {
		val, err := Parser(input)
		if err != nil {
			fmt.Println(err)
		} else {
			output, err := Calculating(val)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Result", output)
			}
		}
	}
}

type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value    interface{}
	priority int8
	next     *Element
}

func (s *Stack) PeekValue() interface{} {
	if s.Len() > 0 {
		return s.top.value
	} else {
		return nil
	}
}

func (s *Stack) PeekPriority() int8 {
	if s.Len() > 0 {
		return s.top.priority
	} else {
		return 0
	}
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) IsEmpty() bool {
	return s.size == 0
}

func (s *Stack) Push(value interface{}, priority int8) {
	if value != ")" {
		s.top = &Element{value, priority, s.top}
		s.size++
	}
}

func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return value
	}
	return nil
}

func (s *Stack) Print() {
	temp := s.top
	fmt.Print("Values stored in stack are: ")
	for temp != nil {
		fmt.Print(temp.value, " ")
		temp = temp.next
	}
	fmt.Println()
}

func (s *Stack) RotateRPN() *Stack {
	tmp := new(Stack)
	for !s.IsEmpty() {
		tmp.Push(s.PeekValue(), s.PeekPriority())
		s.Pop()
	}
	return tmp
}

func (s *Stack) UnarMinus(OperatorRes string) {
	if OperatorRes == "-" {
		s.Push(0.0, 0)
	}
}

func Calculating(s *Stack) (float64, error) {
	FloatStack := new(Stack)
	value := 0.0
	for !s.IsEmpty() {
		switch v := s.Pop().(type) {
		case float64:
			FloatStack.Push(v, 0)
		case string:
			if IsOperator(v) != "" {
				OperRes, err := OperCalc(v, FloatStack.Pop().(float64), FloatStack.Pop().(float64))
				if err != nil {
					return 0, err
				}
				FloatStack.Push(OperRes, 0)
			} else {
				FuncRes, err := FuncCalc(v, FloatStack.Pop().(float64))
				if err != nil {
					return 0, err
				}
				FloatStack.Push(FuncRes, 0)
			}
		default:
			return 0, errors.New("incorrect input data")
		}
	}
	if !FloatStack.IsEmpty() {
		switch v := FloatStack.Pop().(type) {
		case float64:
			value = v
		default:
			return 0, errors.New("incorrect input data")
		}
	}
	return value, nil
}

func OperCalc(OperationRes string, a, b float64) (float64, error) {
	switch OperationRes {
	case "+":
		return b + a, nil
	case "-":
		return b - a, nil
	case "*":
		return b * a, nil
	case "/":
		return b / a, nil
	case "%":
		return math.Mod(b, a), nil
	case "^":
		return math.Pow(b, a), nil
	default:
		return 0, errors.New("wrong operator")
	}
}

func FuncCalc(funcStr string, a float64) (float64, error) {
	switch funcStr {
	case "sin":
		return math.Sin(a), nil
	case "cos":
		return math.Cos(a), nil
	case "tan":
		return math.Tan(a), nil
	case "asin":
		return math.Asin(a), nil
	case "acos":
		return math.Acos(a), nil
	case "atan":
		return math.Atan(a), nil
	case "sqrt":
		return math.Sqrt(a), nil
	case "log":
		return math.Log10(a), nil
	case "ln":
		return math.Log(a), nil
	default:
		return 0, errors.New("wrong name of function")
	}
}

func CheckInput(input []rune) error {
	allowSymbols := "0123456789+-*/()^%acdostiqrnlxg. "
	operators := "+-*/%^"
	PrevSymbol := ""
	bracketCount := 0
	for _, j := range input {
		if !strings.ContainsRune(allowSymbols, j) {
			return errors.New("wrong input data")
		}
		if j == '(' {
			bracketCount++
		} else if j == ')' {
			bracketCount--
		}
		if strings.ContainsRune(operators, j) && strings.Contains(operators, PrevSymbol) {
			return errors.New("duplicate operators")
		}

		PrevSymbol = string(j)

	}
	if bracketCount != 0 {
		return errors.New("incorrect number of brackets")
	}
	return nil
}

func Parser(input []rune) (*Stack, error) {
	st := new(Stack)
	TemporaryStack := new(Stack)
	digit := ""
	funcStr := ""
	for _, j := range input {
		OperatorRes := IsOperator(string(j))
		FunctionRes := IsFunction(j)
		DigitRes := IsDigit(j)
		if DigitRes != "" {
			digit = digit + string(j)
		}
		if DigitRes == "" && len(digit) > 0 {
			a, _ := strconv.ParseFloat(digit, 64)
			TemporaryStack.Push(a, 0)
			digit = ""
		}
		if FunctionRes != "" {
			funcStr = funcStr + string(j)
		}
		if FunctionRes == "" && len(funcStr) > 0 {
			st.Push(funcStr, 5) //highest priority
			funcStr = ""
		}
		if OperatorRes != "" {
			if st.PeekValue() == "(" {
				TemporaryStack.UnarMinus(OperatorRes)
			}
			if OperatorRes != "(" {
				for st.PeekValue() != nil && SeePriority(OperatorRes) <= st.PeekPriority() {
					PopRes := st.Pop().(string)
					if PopRes == "(" {
						break
					}
					TemporaryStack.Push(PopRes, SeePriority(PopRes))
				}
			}
			st.Push(OperatorRes, SeePriority(OperatorRes))
		}
	}
	if digit != "" {
		a, _ := strconv.ParseFloat(digit, 64)
		TemporaryStack.Push(a, 0)
		digit = ""
	}
	for st.PeekValue() != nil {
		PopRes := st.Pop().(string)
		if PopRes != "(" {
			TemporaryStack.Push(PopRes, SeePriority(PopRes))
		}
	}
	OutputStack := TemporaryStack.RotateRPN()
	return OutputStack, nil
}

func SeePriority(input string) int8 {
	switch input {
	case "+":
		return 2
	case "-":
		return 2
	case "*":
		return 3
	case "/":
		return 3
	case "%":
		return 3
	case "^":
		return 4
	case "(":
		return 1
	case ")":
		return 1
	default:
		return 0
	}
}

func IsOperator(input string) string {
	op := "+-*/^%()"
	for _, j := range op {
		if input == string(j) {
			return string(j)
		}
	}
	return ""
}

func IsFunction(input rune) string {
	AllowSymbols := "acdostiqrnlxg "
	for _, j := range AllowSymbols {
		if input == j {
			return string(input)
		}
	}
	return ""
}

func IsDigit(input rune) string {
	Digits := []rune("0123456789.")
	for _, j := range Digits {
		if input == j {
			return string(input)
		}
	}
	return ""
}
