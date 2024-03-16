package main

import "fmt"

type PartialOrder[T any] interface {
	Greater(other T) bool
}

type Number[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64] struct {
	number T
}

func (e *Number[T]) String() string {
	return fmt.Sprintf("%v", e.number)
}

func (a Number[T]) Greater(b Number[T]) bool {
	return a.number > b.number
}

func MakeNumber[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64](n T) Number[T] {
	return Number[T]{
		n,
	}
}

type Node[T PartialOrder[T]] struct {
	Element T
	Left    *Node[T]
	Right   *Node[T]
}

func (n *Node[T]) String() string {
	if n == nil {
		return "Null"
	}

	l := n.Left.String()
	r := n.Right.String()

	return fmt.Sprintf("{Value: %v, Left: %v, Right: %v}", n.Element, l, r)
}

func MakeNode[T PartialOrder[T]](item T) Node[T] {

	return Node[T]{
		item,
		nil,
		nil,
	}

}
func (node *Node[T]) Insert(item T) {
	var toInsert **Node[T]
	if item.Greater(node.Element) {
		toInsert = &node.Right
	} else {
		toInsert = &node.Left
	}
	if *toInsert == nil {
		x := MakeNode(item)
		*toInsert = &x
	} else {
		(**toInsert).Insert(item)
	}
}

func main() {
	tree := MakeNode(MakeNumber(0))

	tree.Insert(MakeNumber(4))

	tree.Insert(MakeNumber(-1))
	tree.Insert(MakeNumber(3))

	fmt.Printf("Tree: %v", &tree)
}
