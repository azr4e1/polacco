package rpn_test

import (
	"github.com/azr4e1/polacco/rpn"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Number interface {
	int | float64
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
			Input:  "-432",
			Output: -432,
		},
		{
			Input:  "005432",
			Output: 5432,
		},
		{
			Input:  "-005432",
			Output: -5432,
		},
		{
			Input:  "      5432",
			Output: 5432,
		},
		{
			Input:  "      -5432",
			Output: -5432,
		},
		{
			Input:  "5432     ",
			Output: 5432,
		},
		{
			Input:  "-5432    ",
			Output: -5432,
		},
		{
			Input:  "0",
			Output: 0,
		},
		{
			Input:  "-0",
			Output: 0,
		},
	}
	testCasesFloat := []TestCase[float64]{
		{
			Input:  "3.14",
			Output: 3.14,
		},
		{
			Input:  "-3.14",
			Output: -3.14,
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
			Input:  "      -3.14",
			Output: -3.14,
		},
		{
			Input:  "-3.14        ",
			Output: -3.14,
		},
		{
			Input:  "0.00",
			Output: 0,
		},
		{
			Input:  "-0.000000000",
			Output: 0,
		},
		{
			Input:  "-3.140000000",
			Output: -3.14,
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
