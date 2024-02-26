package main

import (
	"errors"
	"fmt"
	"strings"
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

func (n Number) String() string {
	return fmt.Sprintf("%v", n.number)
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

func (cn ComplexNumber) String() string {
	if cn.img < 0 {
		return fmt.Sprintf("%v %vi", cn.real, cn.img)
	} else {
		return fmt.Sprintf("%v +%vi", cn.real, cn.img)
	}
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

func makeMatrix(size_x uint, size_y uint) Matrix {
	to_return := Matrix{
		size_x: size_x,
		size_y: size_y,
		value:  make([][]Number, size_x),
	}
	for i := uint(0); i < size_x; i++ {
		to_return.value[i] = make([]Number, size_y)
	}
	return to_return
}

func (m *Matrix) get(x uint, y uint) (MathObject, error) {
	if x >= m.size_x {
		return nil, errors.New("index X out of range")
	}
	if y >= m.size_y {
		return nil, errors.New("index Y out of range")
	}
	return &m.value[x][y], nil
}
func (m *Matrix) set(x uint, y uint, value MathObject) error {
	if x >= m.size_x {
		return errors.New("index X out of range")
	}
	if y >= m.size_y {
		return errors.New("index Y out of range")
	}
	switch n := value.(type) {
	case *Number:
		m.value[x][y] = *n
	default:
		return errors.New("expected number input value")
	}
	return nil
}
func (m *Matrix) get_size_x() uint {
	return m.size_x
}
func (m *Matrix) get_size_y() uint {
	return m.size_y
}

func (m Matrix) String() string {
	return MatrixToString(&m)
}

type ComplexMatrix struct {
	size_x uint
	size_y uint
	value  [][]ComplexNumber
}

func makeComplexMatrix(size_x uint, size_y uint) ComplexMatrix {
	to_return := ComplexMatrix{
		size_x: size_x,
		size_y: size_y,
		value:  make([][]ComplexNumber, size_x),
	}
	for i := uint(0); i < size_x; i++ {
		to_return.value[i] = make([]ComplexNumber, size_y)
	}
	return to_return
}

func (m *ComplexMatrix) get(x uint, y uint) (MathObject, error) {
	if x >= m.size_x {
		return nil, errors.New("index X out of range")
	}
	if y >= m.size_y {
		return nil, errors.New("index Y out of range")
	}
	return &m.value[x][y], nil
}
func (m *ComplexMatrix) set(x uint, y uint, value MathObject) error {
	if x >= m.size_x {
		return errors.New("index X out of range")
	}
	if y >= m.size_y {
		return errors.New("index Y out of range")
	}
	switch n := value.(type) {
	case *ComplexNumber:
		m.value[x][y] = *n
	default:
		return errors.New("expected number input value")
	}
	return nil
}
func (m *ComplexMatrix) get_size_x() uint {
	return m.size_x
}
func (m *ComplexMatrix) get_size_y() uint {
	return m.size_y
}

func (m ComplexMatrix) String() string {
	return MatrixToString(&m)
}

func MatrixToString(m MatrixObject) string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for y := uint(0); y < m.get_size_y(); y++ {
		builder.WriteString("\t")
		for x := uint(0); x < m.get_size_x(); x++ {
			i, e := m.get(x, y)

			if e != nil {
				panic(e)
			}

			builder.WriteString(fmt.Sprintf("%v", i))

			if x+1 == m.get_size_x() {
				builder.WriteString("\n")
			} else {
				builder.WriteString(",\t")
			}
		}
	}

	builder.WriteString("]")

	return builder.String()
}

func main() {

	fmt.Println("################# Number Testing #####################")

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

	fmt.Println("################# Matrix Testing #####################")

	m := makeMatrix(3, 3)
	m.set(2, 0, &Number{100})
	m.set(1, 1, &Number{5})
	fmt.Printf("%v", m)
	m2 := makeComplexMatrix(4, 2)
	m2.set(3, 1, &ComplexNumber{1, -4})
	fmt.Println()
	fmt.Printf("%v", m2)
}
