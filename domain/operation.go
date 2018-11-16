package domain

import (
	"math"
	"reflect"
)

type ParamMap map[string]AbstractInterval

type Operation interface {
	Solve(params ParamMap) Operation
	String() string
	priority() byte
}

type add struct {
	operands []Operation
	result   Operation
}

func (o *add) priority() byte {
	return 1
}

func (o *add) Solve(params ParamMap) Operation {
	if o.result != nil {
		return o.result
	}
	var parametric add
	counted := intervalImpl{}
	for _, op := range o.operands {
		foldedOp := op.Solve(params)
		if reflect.TypeOf(foldedOp) == reflect.TypeOf(counted) {
			counted.left += foldedOp.(intervalImpl).left
			counted.right += foldedOp.(intervalImpl).right
		} else {
			parametric.operands = append(parametric.operands, foldedOp)
		}
	}
	if len(parametric.operands) > 0 {
		if counted.left == 0 && counted.right == 0 {
			o.result = &parametric
			return o.result
		}
		parametric.operands = append(parametric.operands, counted)
		o.result = &parametric
		return o.result
	}
	o.result = counted
	return o.result
}

func (o *add) String() string {
	res := o.Solve(ParamMap{})
	var s string
	if reflect.TypeOf(res) == reflect.TypeOf(&add{}) {
		add, _ := res.(*add)
		for i, op := range add.operands {
			s += op.String()
			if i != len(add.operands)-1 {
				s += " + "
			}
		}
		return s
	}
	return res.String()
}

func Add(operands ...Operation) Operation {
	return &add{
		operands: operands,
	}
}

type sub struct {
	a      Operation
	b      Operation
	result Operation
}

func (o *sub) Solve(params ParamMap) Operation {
	if o.result != nil {
		return o.result
	}
	foldedA := o.a.Solve(params)
	foldedB := o.b.Solve(params)
	if reflect.TypeOf(foldedA) == reflect.TypeOf(intervalImpl{}) &&
		reflect.TypeOf(foldedB) == reflect.TypeOf(intervalImpl{}) {
		res := intervalImpl{
			left:  foldedA.(intervalImpl).left - foldedB.(intervalImpl).left,
			right: foldedA.(intervalImpl).right - foldedB.(intervalImpl).right,
		}
		o.result = res
	} else {
		o.result = &sub{
			a: foldedA,
			b: foldedB,
		}
	}
	return o.result
}

func (o *sub) String() string {
	empty := ParamMap{}
	res := o.Solve(empty)
	var s string
	if reflect.TypeOf(res) == reflect.TypeOf(&sub{}) {
		sub, _ := res.(*sub)
		s += sub.a.Solve(empty).String() + " - "
		if sub.b.priority() <= o.priority() {
			s += "(" + sub.b.Solve(empty).String() + ")"
		} else {
			s += sub.b.Solve(empty).String()
		}
		return s
	}
	return res.String()
}

func (o *sub) priority() byte {
	return 1
}

func Sub(a, b Operation) Operation {
	return &sub{
		a: a,
		b: b,
	}
}

type mul struct {
	operands []Operation
	result   Operation
}

func (o *mul) priority() byte {
	return 2
}

func (o *mul) Solve(params ParamMap) Operation {
	if o.result != nil {
		return o.result
	}
	var parametric mul
	counted := intervalImpl{
		left:  1,
		right: 1,
	}
	for _, op := range o.operands {
		foldedOp := op.Solve(params)
		if reflect.TypeOf(foldedOp) == reflect.TypeOf(counted) {
			factor := foldedOp.(intervalImpl)
			min := math.Min(
				math.Min(
					counted.left*factor.left,
					counted.left*factor.right,
				),
				math.Min(
					counted.right*factor.left,
					counted.right*factor.right,
				),
			)
			max := math.Max(
				math.Max(
					counted.left*factor.left,
					counted.left*factor.right,
				),
				math.Max(
					counted.right*factor.left,
					counted.right*factor.right,
				),
			)
			counted.left = min
			counted.right = max
		} else {
			parametric.operands = append(parametric.operands, foldedOp)
		}
	}
	if len(parametric.operands) > 0 {
		if counted.left == 1 && counted.right == 1 {
			o.result = &parametric
			return o.result
		}
		parametric.operands = append(parametric.operands, counted)
		o.result = &parametric
		return o.result
	}
	o.result = counted
	return o.result
}

func (o *mul) String() string {
	res := o.Solve(ParamMap{})
	var s string
	if reflect.TypeOf(res) == reflect.TypeOf(o) {
		mul, _ := res.(*mul)
		for i, op := range mul.operands {
			if op.priority() < o.priority() {
				s += "(" + op.String() + ")"
			} else {
				s += op.String()
			}
			if i != len(mul.operands)-1 {
				s += "*"
			}
		}
		return s
	}
	return res.String()
}

func Mul(operands ...Operation) Operation {
	return &mul{
		operands: operands,
	}
}

type div struct {
	a      Operation
	b      Operation
	result Operation
}

func (o *div) Solve(params ParamMap) Operation {
	if o.result != nil {
		return o.result
	}
	foldedA := o.a.Solve(params)
	foldedB := o.b.Solve(params)
	if reflect.TypeOf(foldedA) == reflect.TypeOf(intervalImpl{}) &&
		reflect.TypeOf(foldedB) == reflect.TypeOf(intervalImpl{}) {
		dividend := foldedA.(intervalImpl)
		divider := foldedB.(intervalImpl)
		min := math.Min(
			math.Min(
				dividend.left/divider.left,
				dividend.left/divider.right,
			),
			math.Min(
				dividend.right*divider.left,
				dividend.right*divider.right,
			),
		)
		max := math.Max(
			math.Max(
				dividend.left*divider.left,
				dividend.left*divider.right,
			),
			math.Max(
				dividend.right*divider.left,
				dividend.right*divider.right,
			),
		)
		res := intervalImpl{
			left:  min,
			right: max,
		}
		o.result = res
	} else {
		o.result = &div{
			a: foldedA,
			b: foldedB,
		}
	}
	return o.result
}

func (o *div) String() string {
	res := o.Solve(ParamMap{})
	var s string
	if reflect.TypeOf(res) == reflect.TypeOf(o) {
		div, _ := res.(*div)
		if div.a.priority() < o.priority() {
			s += "(" + div.a.String() + ")"
		} else {
			s += div.a.String()
		}
		s += "/"
		if div.b.priority() <= o.priority() {
			s += "(" + div.b.String() + ")"
		} else {
			s += div.b.String()
		}
		return s
	}
	return res.String()
}

func (o *div) priority() byte {
	return 3
}

func Div(a, b Operation) Operation {
	return &div{
		a: a,
		b: b,
	}
}
