package domain

type ring interface {
	mul(multipliers ...Interval) Interval
	add(addednds ...Interval) Interval
}

type operation interface {
	priority() byte
	Solve(varMap VarMap) Interval
	String() string
}

type group interface {
	neutral() intervalImpl
	inversed() group
}

type mul struct {
	k           intervalImpl
	operands    []Interval
	invOperands []Interval
}

func (o mul) neutral() intervalImpl {
	return intervalImpl{1, 1}
}

func (o mul) inversed() group {
	o.k = binaryDiv(o.neutral(), o.k)
	o.operands, o.invOperands = o.invOperands, o.operands
	return o
}

func (o mul) priority() byte {
	return 2
}

func (o mul) Solve(varMap VarMap) Interval {
	var res = mul{
		k: o.k,
	}
	if o.k.left == 0 && o.k.right == 0 {
		return intervalImpl{0, 0}
	}
	for _, operand := range o.operands {
		if i, ok := operand.Solve(varMap).(intervalImpl); ok {
			if i.left == 0 && i.right == 0 {
				return intervalImpl{0, 0}
			}
			res.k = binaryMul(res.k, i)
			continue
		}
		if m, ok := operand.Solve(varMap).(mul); ok {
			if m.k.left == 0 && m.k.right == 0 {
				return intervalImpl{0, 0}
			}
			res.k = binaryMul(m.k, res.k)
			res.operands = append(res.operands, m.operands...)
			res.invOperands = append(res.invOperands, m.invOperands...)
			continue
		}
		res.operands = append(res.operands, operand.Solve(varMap))
	}
	for _, operand := range o.invOperands {
		if i, ok := operand.Solve(varMap).(intervalImpl); ok {
			if i.left == 0 || i.right == 0 {
				panic("Division by zero interval")
			}
			res.k = binaryDiv(res.k, i)
			continue
		}
		if m, ok := operand.Solve(varMap).(mul); ok {
			if m.k.left == 0 || m.k.right == 0 {
				panic("Division by zero interval")
			}
			res.k = binaryDiv(res.k, m.k)
			res.operands = append(res.operands, m.invOperands...)
			res.invOperands = append(res.invOperands, m.operands...)
			continue
		}
		res.invOperands = append(res.invOperands, operand.Solve(varMap))
	}
	if len(res.operands) == 0 && len(res.invOperands) == 0 {
		return res.k
	}
	return res
}

func (o mul) String() string {
	var res string
	if o.k != o.neutral() || len(o.operands) == 0 {
		res += o.k.String()
	}
	for _, operand := range o.operands {
		if len(res) != 0 {
			res += " * "
		}
		if operand.priority() < o.priority() {
			res += "(" + operand.String() + ")"
		} else {
			res += operand.String()
		}
	}
	for _, operand := range o.invOperands {
		res += " / "
		if operand.priority() < o.priority() {
			res += "(" + operand.String() + ")"
		} else {
			res += operand.String()
		}
	}
	return res
}

func (o mul) mul(multipliers ...Interval) Interval {
	o.operands = append(o.operands, multipliers...)
	return o
}

func (o mul) add(addednds ...Interval) Interval {
	res := add{
		m:        add{}.neutral(),
		operands: []Interval{o},
	}
	res.operands = append(res.operands, addednds...)
	return o
}

func (o mul) Add(addednds ...Interval) Interval {
	return o.add(addednds...)
}

func (o mul) Sub(subtrahend Interval) Interval {
	return add{
		m:           add{}.neutral(),
		operands:    []Interval{o},
		invOperands: []Interval{subtrahend},
	}
}

func (o mul) Mul(multipliers ...Interval) Interval {
	return o.mul(multipliers...)
}

func (o mul) Div(divider Interval) Interval {
	o.invOperands = append(o.invOperands, divider)
	return o
}

type add struct {
	m           intervalImpl
	operands    []Interval
	invOperands []Interval
}

func (o add) neutral() intervalImpl {
	return intervalImpl{0, 0}
}

func (o add) inversed() group {
	o.m = binarySub(o.neutral(), o.m)
	o.operands, o.invOperands = o.invOperands, o.operands
	return o
}

func (o add) priority() byte {
	return 1
}

func (o add) Solve(varMap VarMap) Interval {
	var res = add{
		m: o.m,
	}
	for _, operand := range o.operands {
		if i, ok := operand.Solve(varMap).(intervalImpl); ok {
			res.m = binaryAdd(res.m, i)
			continue
		}
		if a, ok := operand.Solve(varMap).(add); ok {
			res.m = binaryAdd(a.m, res.m)
			res.operands = append(res.operands, a.operands...)
			res.invOperands = append(res.invOperands, a.invOperands...)
			continue
		}
		res.operands = append(res.operands, operand.Solve(varMap))
	}
	for _, operand := range o.invOperands {
		if i, ok := operand.Solve(varMap).(intervalImpl); ok {
			res.m = binarySub(res.m, i)
			continue
		}
		if a, ok := operand.Solve(varMap).(add); ok {
			res.m = binarySub(res.m, a.m)
			res.operands = append(res.operands, a.invOperands...)
			res.invOperands = append(res.invOperands, a.operands...)
			continue
		}
		res.invOperands = append(res.invOperands, operand.Solve(varMap))
	}
	if len(res.operands) == 0 && len(res.invOperands) == 0 {
		return res.m
	}
	return res
}

func (o add) String() string {
	var res string
	if o.m != o.neutral() {
		res += o.m.String()
	}
	for _, operand := range o.operands {
		if len(res) != 0 {
			res += " + "
		}
		if operand.priority() < o.priority() {
			res += "(" + operand.String() + ")"
		} else {
			res += operand.String()
		}
	}
	for _, operand := range o.invOperands {
		if len(res) != 0 {
			res += " - "
		}
		if operand.priority() < o.priority() {
			res += "(" + operand.String() + ")"
		} else {
			res += operand.String()
		}
	}
	return res
}

func (o add) mul(multipliers ...Interval) Interval {
	res := mul{
		k:        mul{}.neutral(),
		operands: []Interval{o},
	}
	res.operands = append(res.operands, multipliers...)
	return o
}

func (o add) add(addednds ...Interval) Interval {
	o.operands = append(o.operands, addednds...)
	return o
}

func (o add) Add(addednds ...Interval) Interval {
	return o.add(addednds...)
}

func (o add) Sub(subtrahend Interval) Interval {
	o.invOperands = append(o.invOperands, subtrahend)
	return o
}

func (o add) Mul(multipliers ...Interval) Interval {
	return o.mul(multipliers...)
}

func (o add) Div(divider Interval) Interval {
	return mul{
		k:           mul{}.neutral(),
		operands:    []Interval{o},
		invOperands: []Interval{divider},
	}
}
