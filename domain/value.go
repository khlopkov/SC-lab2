package domain

import "strconv"

type Value interface {
	Value() (float64, bool)
	Add(val Value) Value
	Sub(val Value) Value
	Div(val Value) Value
	Mul(val Value) Value
	Reverse() Value
	String() string
}

type Const float64

func (v Const) Value() (float64, bool) {
	return float64(v), true
}

func (v Const) Add(val Value) Value {
	if valB, ok := val.Value(); ok {
		return Const(valB + float64(v))
	}
	return val.Add(v)
}

func (v Const) Sub(val Value) Value {
	if valB, ok := val.Value(); ok {
		return Const(valB - float64(v))
	}
	return val.Sub(v)
}

func (v Const) Div(val Value) Value {
	if valB, ok := val.Value(); ok {
		return Const(valB / float64(v))
	}
	return val.Reverse().Mul(v)
}

func (v Const) Mul(val Value) Value {
	if valB, ok := val.Value(); ok {
		return Const(valB * float64(v))
	}
	return val.Mul(v)
}

func (v Const) String() string {
	return strconv.FormatFloat(float64(v), 'f', 4, 64)
}

func (v Const) Reverse() Value {
	return Const(1 / float64(v))
}
