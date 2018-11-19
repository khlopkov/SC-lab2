package domain

import (
	"errors"
	"regexp"
	"strconv"
)

var varRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

type Interval interface {
	operation
	ring
	Add(addednds ...Interval) Interval
	Sub(subtrahend Interval) Interval
	Mul(multipliers ...Interval) Interval
	Div(divider Interval) Interval
}

func NewInterval(left float64, right float64) Interval {
	return intervalImpl{
		left:  left,
		right: right,
	}
}

type intervalImpl struct {
	left  float64
	right float64
}

func (i intervalImpl) Solve(varMap VarMap) Interval {
	return i
}

func (i intervalImpl) String() string {
	return "[" + strconv.FormatFloat(i.left, 'f', 4, 64) + ", " + strconv.FormatFloat(i.right, 'f', 4, 64) + "]"
}

func (i intervalImpl) priority() byte {
	return 255
}

func (i intervalImpl) mul(multipliers ...Interval) Interval {
	return mul{
		k:        i,
		operands: multipliers,
	}
}

func (i intervalImpl) add(addends ...Interval) Interval {
	return add{
		m:        i,
		operands: addends,
	}
}

func (i intervalImpl) Add(addednds ...Interval) Interval {
	return i.add(addednds...)
}

func (i intervalImpl) Sub(subtrahend Interval) Interval {
	return add{
		m:           i,
		invOperands: []Interval{subtrahend},
	}
}

func (i intervalImpl) Mul(multipliers ...Interval) Interval {
	return i.mul(multipliers...)
}

func (i intervalImpl) Div(divider Interval) Interval {
	return mul{
		k:           i,
		invOperands: []Interval{divider},
	}
}

type VarMap map[string]Interval

func Var(name string) (Interval, error) {
	if !varRegexp.MatchString(name) {
		return nil, errors.New("bad variable name")
	}
	return variable{
		varName: name,
	}, nil
}

type variable struct {
	varName string
}

func (i variable) Solve(varMap VarMap) Interval {
	if varMap[i.varName] != nil {
		return varMap[i.varName]
	}
	return i
}

func (i variable) String() string {
	return i.varName
}

func (i variable) priority() byte {
	return 255
}

func (i variable) mul(multipliers ...Interval) Interval {
	multipliers = append(multipliers, i)
	return mul{
		k:        mul{}.neutral(),
		operands: multipliers,
	}
}

func (i variable) add(addends ...Interval) Interval {
	addends = append(addends, i)
	return add{
		m:        add{}.neutral(),
		operands: addends,
	}
}

func (i variable) Add(addednds ...Interval) Interval {
	return i.add(addednds...)
}

func (i variable) Sub(subtrahend Interval) Interval {
	return add{
		operands:    []Interval{i},
		invOperands: []Interval{subtrahend},
	}
}

func (i variable) Mul(multipliers ...Interval) Interval {
	return i.mul(multipliers...)
}

func (i variable) Div(divider Interval) Interval {
	return mul{
		k:           mul{}.neutral(),
		operands:    []Interval{i},
		invOperands: []Interval{divider},
	}
}
