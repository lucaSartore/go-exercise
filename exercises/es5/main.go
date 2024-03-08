package main

import (
	"errors"
	"fmt"
)

const (
	ZERO int32 = int32('0')
	NINE int32 = int32('9')
)

func atoi(s string) (int32, error) {
	n := int32(0)
	for _, c := range s {

		n *= 10

		cInt := int32(c) - int32('0')

		if cInt > 9 || cInt < 0 {
			return 0, errors.New("string is not only integers")
		}
		n += cInt
	}
	return n, nil
}

func main() {
	fmt.Println("hello world")

	arr := make([]string, 0)

	arr = append(arr, "10")
	arr = append(arr, "15")
	arr = append(arr, "3")
	arr = append(arr, "354")
	arr = append(arr, "23232")
	arr = append(arr, "2332222")

	for _, v := range arr {
		vInt, err := atoi(v)
		if err != nil {
			fmt.Printf("Got error: %v\n", err)
		} else {
			fmt.Printf("atoi of %v is %v\n", v, vInt)
		}
	}
}
