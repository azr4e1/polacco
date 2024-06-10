package rpn_test

import (
	"github.com/azr4e1/polacco/rpn"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Number interface {
	int | float64
}

const Epsilon = 0.0001

func approxEq(x, y, eps float64) bool {
	diff := x - y
	if diff < 0 {
		diff = -diff
	}

	if diff <= eps {
		return true
	}

	return false
}

func approxEqStack(stack1, stack2 []float64) bool {
	if len(stack1) != len(stack2) {
		return false
	}
	for i := 0; i < len(stack1); i++ {
		if ok := approxEq(stack1[i], stack2[i], Epsilon); !ok {
			return false
		}
	}
	return true
}

func TestStringParser_ParsesANumberCorrectly(t *testing.T) {
	t.Parallel()
	type TestCase[N Number] struct {
		Input  string
		Output N
	}
	testCasesInt := []TestCase[int]{
		{
			Input:  "42532454",
			Output: 42532454,
		},
		{
			Input:  "005432",
			Output: 5432,
		},
		{
			Input:  "      5432",
			Output: 5432,
		},
		{
			Input:  "5432     ",
			Output: 5432,
		},
		{
			Input:  "0",
			Output: 0,
		},
	}
	testCasesFloat := []TestCase[float64]{
		{
			Input:  "3.14",
			Output: 3.14,
		},
		{
			Input:  "	3.14",
			Output: 3.14,
		},
		{
			Input:  "3.14     ",
			Output: 3.14,
		},
		{
			Input:  "0.00",
			Output: 0,
		},
		{
			Input:  "3.140000000",
			Output: 3.14,
		},
	}
	for _, tc := range testCasesInt {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := []float64{float64(tc.Output)}
		got := stack.GetValues()

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}

	for _, tc := range testCasesFloat {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := []float64{float64(tc.Output)}
		got := stack.GetValues()

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParser_ReturnsError(t *testing.T) {
	t.Parallel()
	type TestCase[N Number] struct {
		Input string
	}
	testCasesFloat := []TestCase[float64]{
		{
			Input: "3.14.",
		},
		{
			Input: "	3..14",
		},
		{
			Input: "p 3",
		},
		{
			Input: "0.00_",
		},
		{
			Input: "3.1400000.00",
		},
		{
			Input: "3425345 o",
		},
	}

	for _, tc := range testCasesFloat {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err == nil {
			t.Errorf("want error, got nil for test case '%s'", tc.Input)
		}
	}
}

func TestStringParserSum_AddsTwoElements(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCases := []TestCase{
		{
			Input:  "42532454 2 +",
			Output: []float64{42532456.0},
		},
		{
			Input:  "005432 432 +",
			Output: []float64{5864.0},
		},
		{
			Input:  "      5432 0 +",
			Output: []float64{5432.0},
		},
		{
			Input:  "5432     000542345+",
			Output: []float64{547777},
		},
		{
			Input:  "0 0 +",
			Output: []float64{0},
		},
		{
			Input:  "0.14 5.29 +",
			Output: []float64{5.43},
		},
		{
			Input:  "0 0.0000000 +",
			Output: []float64{0},
		},
		{
			Input:  "0.00000 0.0000000 +",
			Output: []float64{0},
		},
		{
			Input:  "1 2.3 +",
			Output: []float64{3.3},
		},
		{
			Input:  "0 3.14 +",
			Output: []float64{3.14},
		},
		{
			Input:  "3.14 0 +",
			Output: []float64{3.14},
		},
		{
			Input:  "0.0000000001 1000000000000000 +",
			Output: []float64{1000000000000000.0000000001},
		},
		{
			Input:  "23 42		+15    ",
			Output: []float64{65, 15},
		},
	}
	for _, tc := range testCases {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !approxEqStack(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParserDiff_SubtractsTwoElements(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCases := []TestCase{
		{
			Input:  "42532454 2 -",
			Output: []float64{42532452.0},
		},
		{
			Input:  "005432 432 -",
			Output: []float64{5000.0},
		},
		{
			Input:  "      5432 0 -",
			Output: []float64{5432.0},
		},
		{
			Input:  "5432     000542345-",
			Output: []float64{-536913},
		},
		{
			Input:  "0 0 -",
			Output: []float64{0},
		},
		{
			Input:  "0.14 5.29 -",
			Output: []float64{-5.15},
		},
		{
			Input:  "0 0.0000000 -",
			Output: []float64{0},
		},
		{
			Input:  "0.00000 0.0000000 -",
			Output: []float64{0},
		},
		{
			Input:  "1 2.3 -",
			Output: []float64{1.0 - 2.3},
		},
		{
			Input:  "0 3.14 -",
			Output: []float64{-3.14},
		},
		{
			Input:  "3.14 0 -",
			Output: []float64{3.14},
		},
		{
			Input:  "0.0000000001 1000000000000000 -",
			Output: []float64{-999999999999999.9999999999},
		},
		{
			Input:  "23 42		-15   ",
			Output: []float64{-19, 15},
		},
	}
	for _, tc := range testCases {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !approxEqStack(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParserMul_MultipliesTwoElements(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCases := []TestCase{
		{
			Input:  "42532454 2 *",
			Output: []float64{85064908},
		},
		{
			Input:  "005432 432 *",
			Output: []float64{2346624},
		},
		{
			Input:  "      5432 0 *",
			Output: []float64{0},
		},
		{
			Input:  "5432     000542345*",
			Output: []float64{2946018040},
		},
		{
			Input:  "0 0 *",
			Output: []float64{0},
		},
		{
			Input:  "0.14 5.29 *",
			Output: []float64{0.7406},
		},
		{
			Input:  "0 0.0000000 *",
			Output: []float64{0},
		},
		{
			Input:  "0.00000 0.0000000 *",
			Output: []float64{0},
		},
		{
			Input:  "0 1 - 2.3 *",
			Output: []float64{-2.3},
		},
		{
			Input:  "0 3.14 *",
			Output: []float64{0},
		},
		{
			Input:  "3.14 0 *",
			Output: []float64{0},
		},
		{
			Input:  "0 0.0000000001 - 1000000000000000 *",
			Output: []float64{-100000.0},
		},
		{
			Input:  "23 42		*15   ",
			Output: []float64{966, 15},
		},
	}
	for _, tc := range testCases {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !approxEqStack(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParserDiv_DividesTwoElements(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCases := []TestCase{
		{
			Input:  "42532454 2 /",
			Output: []float64{21266227},
		},
		{
			Input:  "005432 432 /",
			Output: []float64{12.574074074074074},
		},
		{
			Input:  "    0  5432  /",
			Output: []float64{0},
		},
		{
			Input:  "5432     000542345/",
			Output: []float64{0.010015764872912999},
		},
		{
			Input:  "0 0.1 -  0.1 /",
			Output: []float64{-1},
		},
		{
			Input:  "0.14 5.29 /",
			Output: []float64{0.026465028355387527},
		},
		{
			Input:  "0 0.00000001 /",
			Output: []float64{0},
		},
		{
			Input:  " 2.3 1/",
			Output: []float64{2.3},
		},
		{
			Input:  "1 3.14 /",
			Output: []float64{0.3184713375796178},
		},
		{
			Input:  "3.14 3.14 /",
			Output: []float64{1},
		},
		{
			Input:  " 1000000000000000 0.0001 /",
			Output: []float64{1e+19},
		},
		{
			Input:  "23 42		/15   ",
			Output: []float64{0.5476190476190477, 15},
		},
	}
	for _, tc := range testCases {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !approxEqStack(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParserPow_PowersTwoElements(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCases := []TestCase{
		{
			Input:  "42532454 2 ^",
			Output: []float64{1809009643262116},
		},
		{
			Input:  "005432 0.01 ^",
			Output: []float64{1.0898070112786702},
		},
		{
			Input:  "      5432 0 ^",
			Output: []float64{1},
		},
		{
			Input:  "0     000542345^",
			Output: []float64{0},
		},
		{
			Input:  "0.14 5.29 ^",
			Output: []float64{3.041006222083085e-05},
		},
		{
			Input:  "0 1 - 2 ^",
			Output: []float64{1},
		},
		{
			Input:  "0 1 - 3 ^",
			Output: []float64{-1},
		},
		{
			Input:  "0 3.14 ^",
			Output: []float64{0},
		},
		{
			Input:  "0 3.14 - 0 ^",
			Output: []float64{1},
		},
		{
			Input:  "2 42		^15   ",
			Output: []float64{4398046511104, 15},
		},
	}
	for _, tc := range testCases {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !approxEqStack(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParserSum_ReturnsError(t *testing.T) {
	t.Parallel()
	input := "23452 +"
	stack := rpn.NewStack()
	err := rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestStringParserDiff_ReturnsError(t *testing.T) {
	t.Parallel()
	input := "23452 -   "
	stack := rpn.NewStack()
	err := rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestStringParserMul_ReturnsError(t *testing.T) {
	t.Parallel()
	input := "23452 *   "
	stack := rpn.NewStack()
	err := rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestStringParserDiv_ReturnsError(t *testing.T) {
	t.Parallel()
	input := "23452 /   "
	stack := rpn.NewStack()
	err := rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error, got nil")
	}

	input = "2345234 0 /"
	stack = rpn.NewStack()
	err = rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error for dividing by zero, got nil")
	}
}

func TestStringParserPow_ReturnsError(t *testing.T) {
	t.Parallel()
	input := "23452 ^   "
	stack := rpn.NewStack()
	err := rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error, got nil")
	}

	input = "0 0 ^"
	stack = rpn.NewStack()
	err = rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error for 0 ^ 0, got nil")
	}

	input = "0 1 - 0.3 ^"
	stack = rpn.NewStack()
	err = rpn.StringParser(stack, input)
	if err == nil {
		t.Error("want error for negative ^ decimal, got nil")
	}
}
