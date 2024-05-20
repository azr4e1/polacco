package rpn_test

import (
	"github.com/azr4e1/polacco/rpn"
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Number interface {
	int | float64
}

const Epsion = 0.0001

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

func TestStringParser_Sum(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCasesInt := []TestCase{
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
	}
	for _, tc := range testCasesInt {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output
		got := stack.GetValues()

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestStringParser_Diff(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input  string
		Output []float64
	}
	testCasesInt := []TestCase{
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
	}
	for _, tc := range testCasesInt {
		stack := rpn.NewStack()
		err := rpn.StringParser(stack, tc.Input)
		if err != nil {
			t.Fatal(err)
		}
		want := tc.Output[0]
		got := stack.GetValues()[0]

		// if !cmp.Equal(want, got) {
		// 	t.Error(cmp.Diff(want, got))
		// }
		if !approxEq(want, got, Epsion) {
			t.Errorf("want %f, got %f", want, got)
		}
	}
}
