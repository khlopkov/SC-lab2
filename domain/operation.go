package domain

import (
	"math"
)

//TODO: implement this
func binaryMul(a intervalImpl, b intervalImpl) intervalImpl {
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
	return intervalImpl{min, max}
}

//TODO: implement this
func binaryDiv(a intervalImpl, b intervalImpl) intervalImpl {
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
	return intervalImpl{min, max}
}

//TODO: implement this
func binaryAdd(a intervalImpl, b intervalImpl) intervalImpl {
	return intervalImpl{a.left + b.left, a.right + b.right}
}

//TODO: implement this
func binarySub(a intervalImpl, b intervalImpl) intervalImpl {
	return intervalImpl{a.left - b.left, a.right - b.right}
}
