package domain

import (
	"math"
)

//TODO: implement this
func binaryMul(a constInterval, b constInterval) constInterval {
	min := math.Min(
		math.Min(
			a.left*b.left,
			a.left*b.right,
		),
		math.Min(
			a.right*b.left,
			a.right*b.right,
		),
	)
	max := math.Max(
		math.Max(
			a.left*b.left,
			a.left*b.right,
		),
		math.Max(
			a.right*b.left,
			a.right*b.right,
		),
	)
	return constInterval{min, max}
}

//TODO: implement this
func binaryDiv(a constInterval, b constInterval) constInterval {
	min := math.Min(
		math.Min(
			a.left/b.left,
			a.left/b.right,
		),
		math.Min(
			a.right/b.left,
			a.right/b.right,
		),
	)
	max := math.Max(
		math.Max(
			a.left/b.left,
			a.left/b.right,
		),
		math.Max(
			a.right/b.left,
			a.right/b.right,
		),
	)
	return constInterval{min, max}
}

//TODO: implement this
func binaryAdd(a constInterval, b constInterval) constInterval {
	return constInterval{a.left + b.left, a.right + b.right}
}

//TODO: implement this
func binarySub(a constInterval, b constInterval) constInterval {
	return constInterval{a.left - b.left, a.right - b.right}
}
