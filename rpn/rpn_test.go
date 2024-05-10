package rpn_test

import (
	"math"
	"testing"

	"github.com/azr4e1/polacco/rpn"
	"github.com/google/go-cmp/cmp"
)

func TestNewStack_CreatesAPopulatedStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, 3, 4, 5}
	stack := rpn.NewStack(want...)
	got := stack.GetValues()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestNewStack_CreatesAnEmptyStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	got := stack.GetValues()

	if len(got) != 0 {
		t.Error("want empty stack, got populated")
	}
}

func TestRPNStackPop_ReturnsLastElement(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack(1, 2, 3)
	want := 3
	got, err := stack.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if float64(want) != got {
		t.Errorf("want %f, got %f", float64(want), got)
	}

	newStack := stack.GetValues()
	if !cmp.Equal([]float64{1, 2}, newStack) {
		t.Error("Stack didn't shrink when popping")
	}
}

func TestRPNStackPop_EmptyStackReturnsError(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	_, err := stack.Pop()

	if err == nil {
		t.Error("wanted error, got nil")
	}
}

func TestRPNStackPush_PushesNewValueToStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, 3, 4}
	stack := rpn.NewStack(1, 2, 3)
	stack.Push(4)

	if !cmp.Equal(want, stack.GetValues()) {
		t.Error("item didn't get pushed")
	}
}

func TestRPNStackAdd_AddsTwoElementsFromTheStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, 7}
	stack := rpn.NewStack(1, 2, 3, 4)
	err := stack.Add()

	if err != nil {
		t.Fatal(err)
	}

	got := stack.GetValues()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRPNStackAdd_ReturnsErrorIfNotEnoughElementsInTheStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	err := stack.Add()

	if err == nil {
		t.Error("want error, got nil for empty stack")
	}

	newStack := rpn.NewStack(1)
	err = newStack.Add()
	if err == nil {
		t.Error("want error, got nil for stack with one item")
	}
}

func TestRPNStackDiff_SubtractTwoElementsFromTheStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, -1}
	stack := rpn.NewStack(1, 2, 3, 4)
	err := stack.Diff()

	if err != nil {
		t.Fatal(err)
	}

	got := stack.GetValues()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRPNStackDiff_ReturnsErrorIfNotEnoughElementsInTheStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	err := stack.Diff()

	if err == nil {
		t.Error("want error, got nil for empty stack")
	}

	newStack := rpn.NewStack(1)
	err = newStack.Diff()
	if err == nil {
		t.Error("want error, got nil for stack with one item")
	}
}

func TestRPNStackDiv_DivideTwoElementsFromTheStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, 3.0 / 4.0}
	stack := rpn.NewStack(1, 2, 3, 4)
	err := stack.Div()

	if err != nil {
		t.Fatal(err)
	}

	got := stack.GetValues()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRPNStackDiv_ReturnsErrorIfNotEnoughElementsInTheStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	err := stack.Div()

	if err == nil {
		t.Error("want error, got nil for empty stack")
	}

	newStack := rpn.NewStack(1)
	err = newStack.Div()
	if err == nil {
		t.Error("want error, got nil for stack with one item")
	}
}

func TestRPNStackMul_MultiplyTwoElementsFromTheStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, 3.0 * 4.0}
	stack := rpn.NewStack(1, 2, 3, 4)
	err := stack.Mul()

	if err != nil {
		t.Fatal(err)
	}

	got := stack.GetValues()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRPNStackMul_ReturnsErrorIfNotEnoughElementsInTheStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	err := stack.Mul()

	if err == nil {
		t.Error("want error, got nil for empty stack")
	}

	newStack := rpn.NewStack(1)
	err = newStack.Mul()
	if err == nil {
		t.Error("want error, got nil for stack with one item")
	}
}

func TestRPNStackPow_PowerTwoElementsFromTheStack(t *testing.T) {
	t.Parallel()
	want := []float64{1, 2, math.Pow(3, 4)}
	stack := rpn.NewStack(1, 2, 3, 4)
	err := stack.Pow()

	if err != nil {
		t.Fatal(err)
	}

	got := stack.GetValues()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRPNStackPow_ReturnsErrorIfNotEnoughElementsInTheStack(t *testing.T) {
	t.Parallel()
	stack := rpn.NewStack()
	err := stack.Pow()

	if err == nil {
		t.Error("want error, got nil for empty stack")
	}

	newStack := rpn.NewStack(1)
	err = newStack.Pow()
	if err == nil {
		t.Error("want error, got nil for stack with one item")
	}
}
