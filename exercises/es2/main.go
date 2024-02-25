package main

import "fmt"

func main() {
	array := [...]int{157, 34, 3, 12, 5, 3, 4, 5, 32, 2, 5, 5, 2, 4, 6, 8, 9, 0, 21, 3, 5, 7, 89, 64, 34, 35, 3333, 4545646, 7, 5, 3, 5}
	merge_sort(array[:])

	fmt.Printf("Sorted Array: %v", array)
}

func merge_sort(array []int) {
	support := make([]int, len(array))
	copy(support, array)
	merge_sort_r(array, support)
}

// when the function return both array and support are sorted
func merge_sort_r(array []int, support []int) {

	n := len(array)

	if n <= 1 {
		return
	}

	array_l := array[0 : n/2]
	array_r := array[n/2:]

	support_l := support[0 : n/2]
	support_r := support[n/2:]

	merge_sort_r(array_l, support_l)
	merge_sort_r(array_r, support_r)

	count := 0
	// as long as i am in the loop both slices are never empty
	for len(support_l) != 0 && len(support_r) != 0 {
		if support_l[0] < support_r[0] {
			array[count] = support_l[0]
			count++
			support_l = support_l[1:]
		} else {
			array[count] = support_r[0]
			count++
			support_r = support_r[1:]
		}
	}
	// copy the remaining stuff
	copy(array[count:], support_l)
	count += len(support_r)
	copy(array[count:], support_r)
	copy(support, array)
}
