package main

import (
	"errors"
	"fmt"
)

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

func (number *Number) sum(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		number.number += v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")
	}
}
func (number *Number) sub(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		number.number -= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")

	}
}
func (number *Number) mul(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		number.number *= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")
	}
}

func (number *Number) div(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		if v.number == 0 {
			return errors.New("division by 0 error")
		}
		number.number /= v.number
		return nil
	default:
		return errors.New("Number can be summed only with number")

	}
}

type ComplexNumber struct {
	real float64
	img  float64
}

func (complex *ComplexNumber) sum(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		complex.real += v.number
	case *ComplexNumber:
		complex.real += v.real
		complex.img += v.img
	default:
		return errors.New("Number can be summed only with number")
	}
	return nil
}
func (complex *ComplexNumber) sub(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		complex.real -= v.number
	case *ComplexNumber:
		complex.real -= v.real
		complex.img -= v.img
	default:
		return errors.New("Number can be summed only with number")

	}
	return nil
}
func (complex *ComplexNumber) mul(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		complex.real *= v.number
		complex.img *= v.number
	case *ComplexNumber:
		new_real := complex.real*v.real + complex.img*v.img
		new_img := complex.img*v.real + complex.real*v.img
		complex.real = new_real
		complex.img = new_img
	default:
		return errors.New("Number can be summed only with number")
	}
	return nil
}

func (complex *ComplexNumber) div(other MathObject) error {
	switch v := other.(type) {
	case *Number:
		if v.number == 0 {
			return errors.New("division by 0 error")
		}
		complex.real /= v.number
		complex.img /= v.number
	case *ComplexNumber:
		if v.img == 0 && complex.real == 0 {
			return errors.New("division by 0 error")
		}
		denominator := v.img*v.img + v.real*v.real
		new_real := (complex.real*v.real + complex.img*v.img) / denominator
		new_img := (complex.img*v.real - complex.real*v.img) / denominator
		complex.real = new_real
		complex.img = new_img
	default:
		return errors.New("Number can be summed only with number")

	}
	return nil
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
	r1 := Number{100}
	r2 := Number{10}
	i1 := ComplexNumber{10, 10}
	i2 := ComplexNumber{2, 5}
	i3 := ComplexNumber{4, -1}

	err := r1.sum(&r2)

	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", r1)

	err = r1.sub(&i1)

	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", r1)

	err = i1.sub(&r2)

	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", i1)

	err = i2.div(&i3)

	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", i2)
}
