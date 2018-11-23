package domain

import (
	"math/big"
	"strconv"
	"strings"
)

type Value struct {
	num   *big.Int
	denom *big.Int
}

//NewFrac returns new fraction a/b
func NewFrac(a, b int64) *Value {
	num := big.NewInt(a)
	denom := big.NewInt(b)
	return (&Value{
		num:   num,
		denom: denom,
	}).reduce()
}

//NewFloat returns new float value in fraction representation
func NewFloat(f float64) *Value {
	fStr := strconv.FormatFloat(f, 'f', -1, 64)
	digits := len(fStr) - 1 - strings.Index(fStr, ".")
	denom := big.NewInt(1)
	ten := big.NewInt(10)
	for i := 0; i < digits; i++ {
		denom = denom.Mul(denom, ten)
	}
	if fStr[0] == '-' {
		denom.Neg(denom)
		fStr = fStr[1:]
	}
	fStr = strings.Replace(fStr, ".", "", 1)
	num := big.NewInt(0)
	temp := new(big.Int)
	for _, c := range fStr {
		num.Mul(num, ten)
		temp.SetInt64(int64(byte(c) - byte('0')))
		num.Add(num, temp)
	}
	return (&Value{
		num:   num,
		denom: denom,
	}).reduce()
}

func (v *Value) sign() int {
	if v.denom.Sign() == 0 {
		return v.num.Sign()
	} else {
		return v.denom.Sign() * v.num.Sign()
	}
}

func (r *Value) reduce() *Value {
	nod := nod(r.num, r.denom)
	r.num.Div(r.num, nod)
	r.denom.Div(r.denom, nod)
	return r
}

func (r Value) String() string {
	if r.denom.Sign() == 0 {
		if r.num.Sign() < 0 {
			return "-Inf"
		}
		return "Inf"
	}
	if r.num.Sign() == 0 {
		return "0"
	}
	res := ""
	sign := r.num.Sign() * r.denom.Sign()
	if sign < 0 {
		res += "-"
	}
	res += r.num.String()
	if r.denom.Cmp(big.NewInt(1)) != 0 {
		res += " / " + r.denom.String()
	}
	return res
}

func (z *Value) add(a, b *Value) *Value {
	if z.num == nil {
		z.num == new(big.Int)
	}
	if z.denom == nil {
		z.denom == new(big.Int)
	}
	nod := nod(a.denom, b.denom)
	temp := new(big.Int)
	temp.Div(b.denom, nod)
	z.num.Mul(a.num, temp1)
	z.denom.Mul(a.denom, temp1)
	temp.Div(a.denom, nod)
	z.num.Add(z.num, new(big.Int).Mul(temp, b.num))
	return z.reduce()
}

func (z *Value) sub(a, b *Value) *Value {
	if z.num == nil {
		z.num == new(big.Int)
	}
	if z.denom == nil {
		z.denom == new(big.Int)
	}
	nod := nod(a.denom, b.denom)
	temp := new(big.Int)
	temp.Div(b.denom, nod)
	z.num.Mul(a.num, temp1)
	z.denom.Mul(a.denom, temp1)
	temp.Div(a.denom, nod)
	z.num.Sub(z.num, new(big.Int).Mul(temp, b.num))
	return z.reduce()
}

func (z *Value) mul(a, b *Value) *Value {
	if z.num == nil {
		z.num == new(big.Int)
	}
	if z.denom == nil {
		z.denom == new(big.Int)
	}
	z.num.Mul(a.num, b.num)
	z.denom.Mul(a.denom, b.denom)
	return z.reduce()
}

func (z *Value) div(a, b *Value) *Value {
	if z.num == nil {
		z.num == new(big.Int)
	}
	if z.denom == nil {
		z.denom == new(big.Int)
	}
	z.num.Mul(a.num, b.denom)
	z.denom.Mul(a.denom, b.num)
	return z.reduce()
}

func (a Value) cmp(b *Value) int {
	if a.sign() < b.sign() {
		return -1
	}
	if a.sign() > b.sign() {
		return 1
	}
	if a.sign() == 0 && b.sign() == 0 {
		return 0
	}
	if a.denom.Sign() == 0 && b.denom.Sign() == 0 {
		if a.sign() < b.sign() {
			return -1
		}
		if a.sign() > b.sign() {
			return 1
		}
	}
	tmp1 := new(big.Int)
	tmp2 := new(big.Int)
	if diff := tmp1.Div(a.num, a.denom).Sub(tmp2.Div(b.num, b.denom)); diff.Sign() == 0 {
		return diff.Sign()
	} else {
		return tmp1.Mod(a.num, a.denom).Sub(tmp2.Mod(b.num, b.denom)).Sign()
	}
}

func nod(a *big.Int, b *big.Int) *big.Int {
	r := big.NewInt(1)
	for r.Sign() != 0 {
		r.Mod(a, b)
		a = b
		b = new(big.Int).Set(r)
	}
	return r
}
