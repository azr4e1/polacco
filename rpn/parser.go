package rpn

import (
	"errors"
	"strconv"
	"strings"
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
		err = errors.New("Unknown operation")
	}

	return err
}

type RPNScanner struct {
	token RPNElement
	exp   string
	pos   int
}

func NewRPNScanner(exp string) *RPNScanner {
	return &RPNScanner{exp: exp}
}

func setNumberToken(s *RPNScanner, token string) {
	conv, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return
	}

	s.token = RPNFloat(conv)
}

func checkOperation(s *RPNScanner, token string) bool {
	switch token {
	case OperationSum, OperationPow, OperationDiff, OperationDiv, OperationMul:
		s.token = RPNOperation(token)
		return true
	default:
		return false
	}

}

func (s *RPNScanner) Scan() bool {
	if s.pos == len(s.exp) {
		return false
	}

	outside := true
	isNum := false
	token := ""
	var i int
	for i = s.pos; i < len(s.exp); i++ {
		s.pos = i
		curr := s.exp[i]

		if outside && unicode.IsDigit(rune(curr)) {
			outside = false
			isNum = true
			token = string(s.exp[i])
			continue
		}
		if !outside && isNum && !(unicode.IsDigit(rune(curr)) || curr == '.') {
			setNumberToken(s, token)
			return true
		}

		if outside && !(strings.Contains(IgnoreCharacters, string(curr)) || unicode.IsDigit(rune(curr))) {
			outside = false
			isNum = false
			token = string(s.exp[i])
			continue
		}
		if !outside && !isNum && (strings.Contains(IgnoreCharacters, string(curr)) || unicode.IsDigit(rune(curr))) {
			if ok := checkOperation(s, token); ok {
				return true
			}
			outside = true
			i--
			continue
		}
		if strings.Contains(IgnoreCharacters, string(curr)) {
			continue
		}

		token += string(curr)
	}
	s.pos = i

	if isNum {
		setNumberToken(s, token)
		return true
	}
	if ok := checkOperation(s, token); ok {
		return true
	}

	return false
}

func (s *RPNScanner) Token() RPNElement {
	return s.token
}

func StringParser(rs *RPNStack, exp string) error {

	scanner := NewRPNScanner(exp)
	for scanner.Scan() {
		token := scanner.Token()
		err := token.Apply(rs)
		if err != nil {
			return err
		}
	}

	return nil
}
