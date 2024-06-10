package rpn

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

const (
	OperationSum  = "+"
	OperationDiff = "-"
	OperationMul  = "*"
	OperationDiv  = "/"
	OperationPow  = "^"
)

const IgnoreCharacters = "\n \t,"

type RPNElement interface {
	Apply(*RPNStack) error
}

type RPNInt int
type RPNFloat float64
type RPNOperation string

func (i RPNInt) Apply(s *RPNStack) error {
	s.Push(float64(i))

	return nil
}

func (f RPNFloat) Apply(s *RPNStack) error {
	s.Push(float64(f))

	return nil
}

func (o RPNOperation) Apply(s *RPNStack) error {
	var err error
	switch o {
	case OperationSum:
		err = s.Add()
	case OperationDiff:
		err = s.Diff()
	case OperationDiv:
		err = s.Div()
	case OperationMul:
		err = s.Mul()
	case OperationPow:
		err = s.Pow()
	default:
		err = errors.New("unknown operation")
	}

	return err
}

type RPNScanner struct {
	token RPNElement
	input string
	pos   int
	err   error
}

func NewRPNScanner(exp string) *RPNScanner {
	return &RPNScanner{input: exp}
}

func isFloatChar(exp rune) bool {
	return unicode.IsDigit(exp) || exp == '.'
}

func (s *RPNScanner) Scan() bool {
	for s.pos < len(s.input) {
		ch := s.input[s.pos]
		switch {
		case unicode.IsDigit(rune(ch)):
			start := s.pos
			periodCounter := 0
			for s.pos < len(s.input) && isFloatChar(rune(s.input[s.pos])) {
				if s.input[s.pos] == '.' {
					periodCounter++
				}
				if periodCounter >= 2 {
					break
				}
				s.pos++
			}
			token := s.input[start:s.pos]
			conv, err := strconv.ParseFloat(token, 64)
			s.token = RPNFloat(conv)
			s.err = err
			return true

		case string(ch) == OperationSum:
			s.pos++
			s.token = RPNOperation(ch)
			return true

		case string(ch) == OperationDiff:
			s.pos++
			s.token = RPNOperation(ch)
			return true

		case string(ch) == OperationMul:
			s.pos++
			s.token = RPNOperation(ch)
			return true

		case string(ch) == OperationDiv:
			s.pos++
			s.token = RPNOperation(ch)
			return true

		case string(ch) == OperationPow:
			s.pos++
			s.token = RPNOperation(ch)
			return true

		case unicode.IsSpace(rune(ch)):
			s.pos++

		default:
			s.pos++
			s.token = nil
			s.err = fmt.Errorf("unexpected character: %s", string(ch))
			return true
		}
	}

	return false
}

func (s *RPNScanner) Token() (RPNElement, error) {
	token, err := s.token, s.err
	s.token, s.err = nil, nil
	return token, err
}

func StringParser(rs *RPNStack, exp string) error {
	scanner := NewRPNScanner(exp)
	for scanner.Scan() {
		token, err := scanner.Token()
		if err != nil {
			return err
		}
		err = token.Apply(rs)
		if err != nil {
			return err
		}
	}

	return nil
}
