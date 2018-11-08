package main

import (
	"fmt"

	"github.com/khlopkov/SC-lab1/domain"
)

func main() {
	var op domain.Operation
	op = domain.Add(
		domain.Div(
			domain.Interval(domain.Const(0), domain.Const(1)),
			domain.Interval(domain.Const(2), domain.Const(5)),
		),
		domain.Mul(
			domain.Interval(domain.Const(0), domain.Const(1)),
			domain.Interval(domain.Const(2), domain.Const(0)),
		),
	)
	fmt.Println(op.Result().String())
}
