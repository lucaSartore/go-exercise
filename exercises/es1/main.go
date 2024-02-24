/*
	Factorial
*/

package main

import (
	"errors"
	"fmt"
)

func main() {
	for i := uint(0); i < 10; i++ {
		test(i)
	}
}

func test(n uint) {
	fmt.Printf("Factorial (iterative) of %v is %v\n", n, factorial_iterative(n))
	fmt.Printf("Factorial (recursive) of %v is %v\n", n, factorial_recursive(n))
	n2 := int(n) - 6
	res, err := factorial_with_error(n2)
	fmt.Printf("Factorial (with error) of %v is %v, %v\n", n2, res, err)
}

func factorial_iterative(n uint) uint {
	to_return := uint(1)
	for i := uint(2); i <= n; i++ {
		to_return *= i
	}
	return uint(to_return)
}

func factorial_recursive(n uint) uint {
	if n == 0 {
		return 1
	}
	return n * factorial_recursive(n-1)
}

func factorial_with_error(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("negative number")
	}
	return int(factorial_iterative(uint(n))), nil
}
