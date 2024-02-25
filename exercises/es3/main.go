package main

import "errors"

type MathObject interface {
	sum(other MathObject) error
	sub(other MathObject) error
	mul(other MathObject) error
	div(other MathObject) error
}

type MatrixObject interface {
	get(x uint, y uint) (MathObject, error)
	set(x uint, y uint, value MathObject) error
	get_size_x() uint
	get_size_y() uint
}

//////////////////////// Real Number /////////////////////////////////////////////

type Number struct {
	number float64
}

func (this Number) sum(other MathObject) error {
	switch v := other.(type) {
	case Number:
		this.number += v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")

	}
}
func (this Number) sub(other MathObject) error {
	switch v := other.(type) {
	case Number:
		this.number -= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")

	}
}
func (this Number) mul(other MathObject) error {
	switch v := other.(type) {
	case Number:
		this.number *= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")
	}
}

func (this Number) div(other MathObject) error {
	switch v := other.(type) {
	case Number:
		if v.number == 0 {
			return errors.New("Division by 0 error")
		}
		this.number /= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")

	}
}

type ComplexNumber struct {
	real float64
	img  float64
}

type Matrix struct {
	size_x uint
	size_y uint
	value  [][]Number
}

type ComplexMatrix struct {
	size_x uint
	size_y uint
	value  [][]ComplexNumber
}

func main() {

}
