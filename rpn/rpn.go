package rpn

import (
	"errors"
	"math"
)

type RPNStack struct {
	values []float64
}

func NewStack(val ...float64) *RPNStack {
	return &RPNStack{values: val}
}

func (r *RPNStack) GetValues() []float64 {
	vals := make([]float64, len(r.values))
	copy(vals, r.values)

	return vals
}

func (r *RPNStack) Pop() (float64, error) {
	length := len(r.values)
	if length < 1 {
		return 0, errors.New("stack is empty")
	}

	var val float64
	val, r.values = r.values[length-1], r.values[:length-1]

	return val, nil
}

func (r *RPNStack) Push(item float64) {
	r.values = append(r.values, item)
}

func (r *RPNStack) Add() error {
	if len(r.values) < 2 {
		return errors.New("not enough elements in the stack")
	}

	var item2, item1 float64
	var err error

	item2, err = r.Pop()
	if err != nil {
		return err
	}

	item1, err = r.Pop()
	if err != nil {
		return err
	}

	sum := item1 + item2
	r.Push(sum)

	return nil
}

func (r *RPNStack) Diff() error {
	if len(r.values) < 2 {
		return errors.New("not enough elements in the stack")
	}

	var item2, item1 float64
	var err error

	item2, err = r.Pop()
	if err != nil {
		return err
	}

	item1, err = r.Pop()
	if err != nil {
		return err
	}

	diff := item1 - item2
	r.Push(diff)

	return nil
}

func (r *RPNStack) Div() error {
	if len(r.values) < 2 {
		return errors.New("not enough elements in the stack")
	}

	var item2, item1 float64
	var err error

	item2, err = r.Pop()
	if err != nil {
		return err
	}
	if item2 == 0 {
		return errors.New("cannot divide by 0")
	}

	item1, err = r.Pop()
	if err != nil {
		return err
	}

	div := item1 / item2
	r.Push(div)

	return nil
}

func (r *RPNStack) Mul() error {
	if len(r.values) < 2 {
		return errors.New("not enough elements in the stack")
	}

	var item2, item1 float64
	var err error

	item2, err = r.Pop()
	if err != nil {
		return err
	}

	item1, err = r.Pop()
	if err != nil {
		return err
	}

	mul := item1 * item2
	r.Push(mul)

	return nil
}

func (r *RPNStack) Pow() error {
	if len(r.values) < 2 {
		return errors.New("not enough elements in the stack")
	}

	var item2, item1 float64
	var err error

	item2, err = r.Pop()
	if err != nil {
		return err
	}

	item1, err = r.Pop()
	if err != nil {
		return err
	}

	if item1 == 0 && item2 == 0 {
		return errors.New("cannot raise 0 to the power of 0")
	}

	if item1 < 0 && item2-math.Trunc(item2) != 0 {
		return errors.New("cannot raise a negative number to a fractional exponent")
	}

	pow := math.Pow(item1, item2)
	r.Push(pow)

	return nil
}
