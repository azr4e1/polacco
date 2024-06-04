package shell_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/azr4e1/polacco/shell"
	"github.com/google/go-cmp/cmp"
)

func TestShellRun_RunsListCorrectly(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 + +  \nls\n"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "> > [6]\n> "
	session.Run()
	got := output.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRun_RunsPopCorrectly(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \npop\n"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "> > 3\n> "
	session.Run()
	got := output.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRun_PopReturnsErrorCorrectly(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "pop\n"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "error: stack is empty\n"
	session.Run()
	got := error.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRun_ResetRunsCorrectly(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \nreset\nls\n"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "> > > []\n> "
	session.Run()
	got := output.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRun_ReturnsErrorForParsingErrors(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 + + +"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "error: not enough elements in the stack\n"
	session.Run()
	got := error.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}

	input = new(bytes.Buffer)
	output = new(bytes.Buffer)
	error = new(bytes.Buffer)
	session, err = shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr = "3 0 /"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want = "error: cannot divide by 0\n"
	session.Run()
	got = error.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}

	input = new(bytes.Buffer)
	output = new(bytes.Buffer)
	error = new(bytes.Buffer)
	session, err = shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetPrompt("> "),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr = "0 0 ^"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want = "error: cannot raise 0 to the power of 0\n"
	session.Run()
	got = error.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}

	input = new(bytes.Buffer)
	output = new(bytes.Buffer)
	error = new(bytes.Buffer)
	session, err = shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr = "0 1 - 1.5 ^"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want = "error: cannot raise a negative number to a fractional exponent\n"
	session.Run()
	got = error.String()

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRun_StoresHistoryCorrectly(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \npop\np\nh\nl  \nls\n^\n-\n-\n	\n-\n	\n"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := []string{"3 1 2 +  ", "pop", "p", "h", "l  ", "ls", "^", "-"}
	session.Run()
	got := session.GetHistory()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestShellRun_MaintainsTheCorrectHistoryLength(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	maxLength := 100
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
		shell.SetMaxHistorySize(maxLength),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "pop"
	for i := 0; i < maxLength; i++ {
		inputStr += fmt.Sprintf("\nls%d", i)
	}
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := []string{}
	for i := 0; i < maxLength; i++ {
		want = append(want, fmt.Sprintf("ls%d", i))
	}

	session.Run()
	got := session.GetHistory()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestShellRunGetPrevHistoryElement_ReturnsTheCorrectElements(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \npop\np\nh\nl  \nls\n^\n-"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "-"
	session.Run()
	got, err := session.GetPrevHistoryElement()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRunGetNextHistoryElement_ReturnsTheCorrectElements(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \npop\np\nh\nl  \nls\n^\n-"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	want := "-"
	session.Run()
	_, err = session.GetPrevHistoryElement()
	if err != nil {
		t.Fatal(err)
	}
	_, err = session.GetPrevHistoryElement()
	if err != nil {
		t.Fatal(err)
	}
	got, err := session.GetNextHistoryElement()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestShellRunGetPrevHistoryElement_ReturnsError(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  "
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	session.Run()
	_, err = session.GetPrevHistoryElement()
	_, err = session.GetPrevHistoryElement()
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestShellRunGetNextHistoryElement_ReturnsError(t *testing.T) {
	t.Parallel()
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)
	error := new(bytes.Buffer)
	session, err := shell.NewSession(
		shell.SetStdin(input),
		shell.SetStdout(output),
		shell.SetStderr(error),
	)
	if err != nil {
		t.Error(err)
	}
	inputStr := "3 1 2 +  \npop\np\nh\nl  \nls\n^\n-"
	_, err = input.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}
	session.Run()
	_, err = session.GetNextHistoryElement()
	if err == nil {
		t.Error("want error, got nil")
	}
}
