package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/azr4e1/polacco/rpn"
)

var Help = `
h: print this help
q: quit
p: pop and show last element of stack
r: reset stack
l: show stack

Supported operations:
	+: sum
	-: diff
	/: div
	*: mul
	^: pow
`

type session struct {
	input          io.Reader
	output         io.Writer
	error          io.Writer
	stack          *rpn.RPNStack
	history        []string
	maxHistorySize int
	historyPointer int
}

type option func(*session) error

func NewSession(opts ...option) (*session, error) {
	stack := rpn.NewStack()
	s := &session{
		input:          os.Stdin,
		output:         os.Stdout,
		error:          os.Stderr,
		stack:          stack,
		maxHistorySize: 50,
	}

	for _, o := range opts {
		err := o(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func SetStdin(stdin io.Reader) option {
	return func(s *session) error {
		if stdin == nil {
			return errors.New("stdin is nil")
		}

		s.input = stdin
		return nil
	}
}

func SetStdout(stdout io.Writer) option {
	return func(s *session) error {
		if stdout == nil {
			return errors.New("stdout is nil")
		}

		s.output = stdout
		return nil
	}
}

func SetStderr(stderr io.Writer) option {
	return func(s *session) error {
		if stderr == nil {
			return errors.New("stderr is nil")
		}

		s.error = stderr
		return nil
	}
}

func SetStack(vals ...float64) option {
	return func(s *session) error {
		stack := rpn.NewStack(vals...)
		s.stack = stack

		return nil
	}
}

func SetMaxHistorySize(maxHistorySize int) option {
	return func(s *session) error {
		if maxHistorySize < 0 {
			return errors.New("cannot set negative history size")
		}

		s.maxHistorySize = maxHistorySize
		return nil
	}
}

func (s *session) Exec(expr string) {
	cleanExpr := strings.ToLower(strings.TrimSpace(expr))
	switch cleanExpr {
	case "h", "he", "hel", "help":
		s.Help()
	case "l", "ls", "li", "lis", "list":
		s.List()
	case "p", "po", "pop":
		s.Pop()
	case "r", "re", "res", "rese", "reset":
		s.Reset()
	case "q", "qu", "qui", "quit":
		os.Exit(0)
	default:
		s.Parse(cleanExpr)
	}
}

func (s *session) Help() {
	fmt.Fprintln(s.output, Help)
}

func (s *session) List() {
	fmt.Fprintf(s.output, "%v\n", s.stack.GetValues())
}

func (s *session) Pop() {
	val, err := s.stack.Pop()
	if err != nil {
		fmt.Fprintln(s.error, "error:", err)
		return
	}
	fmt.Fprintln(s.output, val)
}

func (s *session) Reset() {
	stack := rpn.NewStack()
	s.stack = stack
}

func (s *session) Parse(expr string) {
	err := rpn.StringParser(s.stack, expr)
	if err != nil {
		fmt.Fprintln(s.error, "error:", err)
	}
}

func (s *session) GetHistory() []string {
	return s.history
}

func (s *session) updateHistory(expr string) {
	// ignore empty
	if strings.TrimSpace(expr) == "" {
		return
	}
	// ignore repetition
	if length := len(s.history); length > 0 && s.history[length-1] == expr {
		return
	}
	s.history = append(s.history, expr)
	s.history = s.history[max(0, len(s.history)-s.maxHistorySize):len(s.history)]
	s.historyPointer = len(s.history) // reset pointer to last element
}

func (s *session) GetPrevHistoryElement() (string, error) {
	if s.historyPointer <= 0 {
		s.historyPointer = 0
		return "", errors.New("earliest point in history")
	}
	s.historyPointer--
	return s.history[s.historyPointer], nil
}

func (s *session) GetNextHistoryElement() (string, error) {
	if s.historyPointer >= len(s.history)-1 {
		s.historyPointer = len(s.history) - 1
		return "", errors.New("latest point in history")
	}
	s.historyPointer++
	return s.history[s.historyPointer], nil
}

func (s *session) Run() {
	fmt.Fprintf(s.output, "> ")
	scanner := bufio.NewScanner(s.input)
	for scanner.Scan() {
		expr := scanner.Text()
		s.updateHistory(expr)

		s.Exec(expr)
		fmt.Fprintf(s.output, "> ")
	}
}
