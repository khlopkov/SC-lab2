package domain

type ring interface {
	mul(multiplier operation) operation
	add(addednd operation) operation
}

type operation interface {
	ring
	priority() byte
	Solve(varMap VarMap) operation
	String() string
}

type group interface {
	neutral() constInterval
	inversed() group
	op() operation
}

type mul struct {
	k           constInterval
	operands    []operation
	invOperands []operation
}

func (o mul) neutral() constInterval {
	return constInterval{1, 1}
}

func (o mul) inversed() group {
	o.k = binaryDiv(o.neutral(), o.k)
	o.operands, o.invOperands = o.invOperands, o.operands
	return o
}

func (o mul) op() operation {
	return o
}

func (o mul) priority() byte {
	return 2
}

func (o mul) Solve(varMap VarMap) operation {
	var res = mul{
		k: o.k,
	}
	if o.k.left == 0 && o.k.right == 0 {
		return constInterval{0, 0}
	}
	for _, operand := range o.operands {
		if i, ok := operand.Solve(varMap).(constInterval); ok {
			if i.left == 0 && i.right == 0 {
				return constInterval{0, 0}
			}
			res.k = binaryMul(res.k, i)
			continue
		}
		if m, ok := operand.Solve(varMap).(mul); ok {
			if m.k.left == 0 && m.k.right == 0 {
				return constInterval{0, 0}
			}
			res.k = binaryMul(m.k, res.k)
			res.operands = append(res.operands, m.operands...)
			res.invOperands = append(res.invOperands, m.invOperands...)
			continue
		}
		res.operands = append(res.operands, operand.Solve(varMap))
	}
	for _, operand := range o.invOperands {
		if i, ok := operand.Solve(varMap).(constInterval); ok {
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

func (o mul) mul(multiplier operation) operation {
	o.operands = append(o.operands, multiplier)
	return o
}

func (o mul) add(addednd operation) operation {
	res := add{
		m:        add{}.neutral(),
		operands: []operation{o},
	}
	res.operands = append(res.operands, addednd)
	return res
}

type add struct {
	m           constInterval
	operands    []operation
	invOperands []operation
}

func (o add) neutral() constInterval {
	return constInterval{0, 0}
}

func (o add) inversed() group {
	o.m = binarySub(o.neutral(), o.m)
	o.operands, o.invOperands = o.invOperands, o.operands
	return o
}

func (o add) op() operation {
	return o
}

func (o add) priority() byte {
	return 1
}

func (o add) Solve(varMap VarMap) operation {
	var res = add{
		m: o.m,
	}
	for _, operand := range o.operands {
		if i, ok := operand.Solve(varMap).(constInterval); ok {
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
		if i, ok := operand.Solve(varMap).(constInterval); ok {
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

func (o add) mul(multiplier operation) operation {
	res := mul{
		k:        mul{}.neutral(),
		operands: []operation{o},
	}
	res.operands = append(res.operands, multiplier)
	return o
}

func (o add) add(addednd operation) operation {
	o.operands = append(o.operands, addednd)
	return o
}
