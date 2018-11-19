package main

import (
	"fmt"

	"github.com/khlopkov/SC-lab2/domain"
)

func main() {
	variable, err := domain.Var("x")
	if err != nil {
		panic(err)
	}

	op := variable.Add(
		variable.Div(
			domain.NewInterval(5, 5).Add(
				variable,
			),
		).Mul(
			domain.NewInterval(1, 1),
		).Div(
			variable,
		),
	).Sub(domain.NewInterval(1, 1))
	fmt.Println(op.Solve(domain.VarMap{"u": domain.NewInterval(1, 1)}).String())
}
