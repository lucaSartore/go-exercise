package main

type MathObject interface {
	Sum(a MathObject, b MathObject) (MathObject, error)
	Sub(a MathObject, b MathObject) (MathObject, error)
	Mul(a MathObject, b MathObject) (MathObject, error)
	Div(a MathObject, b MathObject) (MathObject, error)
}

func main() {

}
