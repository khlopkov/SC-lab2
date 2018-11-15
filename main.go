package main

import (
	"fmt"

	"github.com/khlopkov/SC-lab1/domain"
)

func main() {
	variable, err := domain.ParametricInterval("x")
	if err != nil {
		panic(err)
	}

	op := domain.Add(
		domain.Interval(1, 4),
		domain.Div(
			variable,
			domain.Add(
				variable,
				domain.Interval(2, -5),
			),
		),
	)
	fmt.Println(op.Solve(domain.ParamMap{"u": domain.Interval(-3, 4)}).String())
}
