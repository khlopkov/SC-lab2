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
			domain.NewInterval(domain.NewFrac(5, 1), domain.NewFrac(6, 1)).Add(
				variable,
			),
		).Mul(
			domain.NewInterval(domain.NewFrac(5, 3), domain.NewFrac(1, 2)),
		).Div(
			variable,
		),
	).Sub(domain.NewInterval(domain.NewFrac(1, 1), domain.NewFrac(1, 1)))
	fmt.Println(op.Solve(domain.VarMap{"u": domain.NewInterval(domain.NewFrac(1, 1), domain.NewFrac(1, 1))}).String())
}
