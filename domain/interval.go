/*
Package domain describes intervals and base operation on them:
<ul>
	<li> Add </li>
	<li> Sub </li>
	<li> Mul </li>
	<li> Div </li>
</ul>
*/
package domain

import (
	"errors"
	"regexp"
)

var varRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

//Interval structure describing interval with operations on it
type Interval struct {
	op operation
}

//NewInterval creates new interval with const values of left and right bounds of interval
func NewInterval(left *Value, right *Value) Interval {
	return Interval{
		op: constInterval{
			left:  left,
			right: right,
		},
	}
}

//String returns string representation of interval
func (i Interval) String() string {
	return i.op.String()
}

//Solve folds and solve interval with variable values passed in VarMap. Returns folded and solved interval
func (i Interval) Solve(varMap VarMap) Interval {
	i.op = i.op.Solve(varMap)
	return i
}

type constInterval struct {
	left  *Value
	right *Value
}

func (i constInterval) Solve(varMap VarMap) operation {
	return i
}

func (i constInterval) String() string {
	return "[" + i.left.String() + ", " + i.right.String() + "]"
}

func (i constInterval) priority() byte {
	return 255
}

func (i constInterval) mul(multiplier operation) operation {
	return mul{
		k:        i,
		operands: []operation{multiplier},
	}
}

func (i constInterval) add(addend operation) operation {
	return add{
		m:        i,
		operands: []operation{addend},
	}
}

func (a constInterval) addConst(b constInterval) constInterval {
	res := constInterval{
		left:  new(Value),
		right: new(Value),
	}
	res.left.add(a.left, b.left)
	if res.left.cmp(res.right) > 0 {
		res.right, res.left = res.left, res.right
	}
	return res
}

func (a constInterval) subConst(b constInterval) constInterval {
	res := constInterval{
		left:  new(Value),
		right: new(Value),
	}
	res.left.sub(a.left, b.left)
	if res.left.cmp(res.right) > 0 {
		res.right, res.left = res.left, res.right
	}
	return res
}

func (a constInterval) mulConst(b constInterval) constInterval {
	temp := new(Value)
	min := new(Value)
	min.mul(a.left, b.left)
	max := new(Value)
	max.mul(a.left, b.left)

	temp.mul(a.left, b.right)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}

	temp.mul(a.right, b.left)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}

	temp.mul(a.right, b.right)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}
	return constInterval{min, max}
}

func (a constInterval) divConst(b constInterval) constInterval {
	temp := new(Value)
	min := new(Value)
	min.div(a.left, b.left)
	max := new(Value)
	max.div(a.left, b.left)

	temp.div(a.left, b.right)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}

	temp.div(a.right, b.left)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}

	temp.div(a.right, b.right)
	if temp.cmp(min) < 0 {
		min = temp
	}
	if temp.cmp(max) > 0 {
		max = temp
	}
	return constInterval{min, max}
}

//Add returns result of addition current interval and passed addends
//[a, b].Add([a1,b1], [a2, b2], ..., [an, bn]) = [a, b] + [a1, b1] + [a2, b2] + ... + [an, bn]
func (i Interval) Add(addednds ...Interval) Interval {
	op := i.op
	for _, a := range addednds {
		op = op.add(a.op)
	}
	return Interval{
		op: op,
	}
}

//Sub returns result of subtraction subtrahend from current interval
//[a, b].Sub([c, d]) = [a, b] - [c, d]
func (i Interval) Sub(subtrahend Interval) Interval {
	op := add{
		m:        add{}.neutral(),
		operands: []operation{subtrahend.op},
	}
	i.op = i.op.add(op.inversed().op())
	return i
}

//Mul returns result of multiplying current interval on passed multipliers
//[a, b].Mul([a1,b1], [a2, b2], ..., [an, bn]) = [a, b] * [a1, b1] * [a2, b2] * ... * [an, bn]
func (i Interval) Mul(multipliers ...Interval) Interval {
	i.op = mul{
		k:        mul{}.neutral(),
		operands: []operation{i.op},
	}
	for _, m := range multipliers {
		i.op = i.op.mul(m.op)
	}
	return i
}

//Div returns result of division current interval on divider
//[a, b].Sub([c, d]) = [a, b] / [c, d]
func (i Interval) Div(divider Interval) Interval {
	op := mul{
		k:        mul{}.neutral(),
		operands: []operation{divider.op},
	}
	i.op = i.op.mul(op.inversed().op())
	return i
}

//VarMap type describing variable values in format {"varName", Interval}
type VarMap map[string]Interval

//Var creates new variable interval with passed name.
//variable name should contain only letters and digits and start from letter, else creation will return error
func Var(name string) (Interval, error) {
	if !varRegexp.MatchString(name) {
		return Interval{}, errors.New("bad variable name")
	}
	return Interval{
		op: variable{
			varName: name,
		},
	}, nil
}

type variable struct {
	varName string
}

func (i variable) Solve(varMap VarMap) operation {
	if varMap[i.varName].op != nil {
		return varMap[i.varName].op
	}
	return i
}

func (i variable) String() string {
	return i.varName
}

func (i variable) priority() byte {
	return 255
}

func (i variable) mul(multiplier operation) operation {
	return mul{
		k:        mul{}.neutral(),
		operands: []operation{i, multiplier},
	}
}

func (i variable) add(addend operation) operation {
	return add{
		m:        add{}.neutral(),
		operands: []operation{addend, i},
	}
}
