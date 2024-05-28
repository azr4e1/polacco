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

var Help = `Ok homeboy

`

type session struct {
	input          io.Reader
	output         io.Writer
	error          io.Writer
	stack          *rpn.RPNStack
	history        []string
	maxHistorySize int
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
		fmt.Fprintln(s.output, Help)
	case "l", "ls", "li", "lis", "list":
		fmt.Fprintf(s.output, "%v\n", s.stack.GetValues())
	case "p", "po", "pop":
		val, err := s.stack.Pop()
		if err != nil {
			fmt.Fprintln(s.error, "error:", err)
			return
		}
		fmt.Fprintln(s.output, val)
	case "r", "re", "res", "rese", "reset":
		stack := rpn.NewStack()
		s.stack = stack
	case "q", "qu", "qui", "quit":
		os.Exit(0)
	default:
		err := rpn.StringParser(s.stack, cleanExpr)
		if err != nil {
			fmt.Fprintln(s.error, "error:", err)
		}
	}
}

func (s *session) Run() {
	fmt.Fprintf(s.output, "> ")
	scanner := bufio.NewScanner(s.input)
	for scanner.Scan() {
		expr := scanner.Text()
		s.history = append(s.history, expr)
		s.history = s.history[:min(len(s.history), s.maxHistorySize)]

		s.Exec(expr)
		fmt.Fprintf(s.output, "> ")
	}
}

func (s *session) GetHistory() []string {
	return s.history
}
