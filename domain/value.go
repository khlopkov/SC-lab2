package domain

import (
	"math/big"
	"strconv"
	"strings"
)

//Inf returns infinite Value
func Inf() *Value {
	return &Value{
		num:   big.NewInt(1),
		denom: big.NewInt(0),
	}
}

//NegInf returns negative infinite value
func NegInf() *Value {
	return &Value{
		num:   big.NewInt(-1),
		denom: big.NewInt(0),
	}
}

//NaN returns NaN value.
//WARNING: comparing two NaN values will lead to panic
//PANIC NaN.cmp(NaN())
func NaN() *Value {
	return &Value{
		num:   big.NewInt(0),
		denom: big.NewInt(0),
	}
}

//Zero returns zero value
func Zero() *Value {
	return &Value{
		num:   big.NewInt(0),
		denom: big.NewInt(1),
	}
}

//Zero returns zero value
func One() *Value {
	return &Value{
		num:   big.NewInt(1),
		denom: big.NewInt(1),
	}
}

//Value type for representation numeric value with fraction
type Value struct {
	num   *big.Int
	denom *big.Int
}

//NewFrac returns new fraction a/b
func NewFrac(a, b int64) *Value {
	if b == 0 {
		if a > 0 {
			return Inf()
		}
		if a < 0 {
			return NegInf()
		}
		return NaN()
	}
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

//NewInt returns new int value in fraction representation
func NewInt(i int64) *Value {
	return (&Value{
		num:   big.NewInt(i),
		denom: big.NewInt(1),
	}).reduce()
}

func (v *Value) sign() int {
	v.checkNil()

	if v.denom.Sign() == 0 {
		return v.num.Sign()
	}
	return v.denom.Sign() * v.num.Sign()

}

func (r *Value) reduce() *Value {
	r.checkNil()

	nod := nod(r.num, r.denom)
	r.num.Div(r.num, nod)
	r.denom.Div(r.denom, nod)
	if r.num.Sign() < 0 && r.denom.Sign() < 0 {
		r.num.Neg(r.num)
		r.denom.Neg(r.denom)
	}
	return r
}

//String converts value to string
func (r Value) String() string {
	r.checkNil()

	if r.denom.Sign() == 0 && r.num.Sign() == 0 {
		return "NaN"
	}
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
	temp := new(big.Int)
	if r.num.Sign() < 0 {
		res += temp.Neg(r.num).String()
	} else {
		res += r.num.String()
	}
	if r.denom.Cmp(big.NewInt(1)) != 0 {
		res += " / "
		if r.denom.Sign() < 0 {
			res += new(big.Int).Neg(r.denom).String()
		} else {
			res += r.denom.String()
		}
	}
	return res
}

//add adds to values
//z = a + b
//Some specific operations
//Inf() + Inf() == Inf()
//NegInf() + NegInf() == NegInf()
//Inf() + x == Inf() for every finite x
//NegInf() + x == NegInf() for every finite x
//Inf() + NegInf() == NaN()
//NaN() + x == NaN() for every x
//NaN() + NaN() == NaN() for every x
func (z *Value) add(a, b *Value) *Value {
	z.checkNil()
	a.checkNil()
	b.checkNil()

	if a.num.Sign() == 0 && a.denom.Sign() == 0 ||
		b.num.Sign() == 0 && b.denom.Sign() == 0 ||
		a.cmp(Inf()) == 0 && b.cmp(NegInf()) == 0 ||
		b.cmp(Inf()) == 0 && a.cmp(NegInf()) == 0 {
		z = NaN()
		return z
	}
	if a.cmp(Inf()) == 0 || b.cmp(Inf()) == 0 {
		z = Inf()
		return z
	}
	if a.cmp(NegInf()) == 0 || b.cmp(NegInf()) == 0 {
		z = NegInf()
		return z
	}

	nod := nod(a.denom, b.denom)
	temp := new(big.Int)
	temp.Div(b.denom, nod)
	z.num.Mul(a.num, temp)
	z.denom.Mul(a.denom, temp)
	temp.Div(a.denom, nod)
	z.num.Add(z.num, new(big.Int).Mul(temp, b.num))
	return z.reduce()
}

//sub subs to values
//z = a - b
//Some specific operations
//Inf() - Inf() == NaN()
//NegInf() - NegInf() == NaN()
//Inf() - x == Inf() for every finite x
//NegInf() - x == NegInf() for every finite x
//Inf() - NegInf() == NaN()
//NaN() - x == NaN() for every x
//NaN() - NaN() == NaN() for every x
func (z *Value) sub(a, b *Value) *Value {
	z.checkNil()
	a.checkNil()
	b.checkNil()

	if a.num.Sign() == 0 && a.denom.Sign() == 0 ||
		b.num.Sign() == 0 && b.denom.Sign() == 0 ||
		a.cmp(Inf()) == 0 && b.cmp(Inf()) == 0 ||
		b.cmp(NegInf()) == 0 && a.cmp(NegInf()) == 0 {
		z = NaN()
		return z
	}
	if a.cmp(Inf()) == 0 {
		z = Inf()
		return z
	}
	if a.cmp(NegInf()) == 0 {
		z = NegInf()
		return z
	}
	if b.cmp(NegInf()) == 0 {
		z = Inf()
		return z
	}
	if b.cmp(Inf()) == 0 {
		z = NegInf()
		return z
	}

	nod := nod(a.denom, b.denom)
	temp := new(big.Int)
	temp.Div(b.denom, nod)
	z.num.Mul(a.num, temp)
	z.denom.Mul(a.denom, temp)
	temp.Div(a.denom, nod)
	z.num.Sub(z.num, new(big.Int).Mul(temp, b.num))
	return z.reduce()
}

//mul multiplies to values
//z = a * b
//Some specific operations
//Inf() * x == Inf() for every positive x
//Inf() * x == NegInf() for every negative x
//NegInf() * x == NegInf() for every positive x
//NegInf() * x == Inf() for every negative x
//Inf() * NegInf() == NaN()
//NaN() * x == NaN() for every x
//NaN() * NaN() == NaN() for every x
func (z *Value) mul(a, b *Value) *Value {
	z.checkNil()
	a.checkNil()
	b.checkNil()

	if a.num.Sign() == 0 && a.denom.Sign() == 0 ||
		b.num.Sign() == 0 && b.denom.Sign() == 0 {
		z = NaN()
		return z
	}
	if a.cmp(Inf()) == 0 || b.cmp(Inf()) == 0 ||
		a.cmp(NegInf()) == 0 || b.cmp(NegInf()) == 0 {
		if a.sign()*b.sign() > 0 {
			z = Inf()
			return z
		}
		z = NegInf()
		return z
	}

	z.num.Mul(a.num, b.num)
	z.denom.Mul(a.denom, b.denom)
	return z.reduce()
}

//div divides to values
//z = a / b
//Some specific operations
//Inf() / Inf() == NaN()
//NegInf() / NegInf() == NaN()
//Inf() / x == Inf() for every positive x
//Inf() / x == NegInf() for every negative x
//NegInf() / x == NegInf() for every positive x
//NegInf() / x == Inf() for every negative x
//Inf() / NegInf() == NaN()
//NaN() / x == NaN() for every x
//NaN() / NaN() == NaN() for every x
func (z *Value) div(a, b *Value) *Value {
	z.checkNil()
	a.checkNil()
	b.checkNil()

	if a.num.Sign() == 0 && a.denom.Sign() == 0 ||
		b.num.Sign() == 0 && b.denom.Sign() == 0 ||
		a.cmp(Inf()) == 0 && b.cmp(NegInf()) == 0 ||
		b.cmp(Inf()) == 0 && a.cmp(NegInf()) == 0 ||
		a.cmp(Inf()) == 0 && b.cmp(Inf()) == 0 ||
		b.cmp(NegInf()) == 0 && a.cmp(NegInf()) == 0 {
		z = NaN()
		return z
	}
	if a.cmp(Inf()) == 0 || a.cmp(NegInf()) == 0 {
		if a.sign()*b.sign() > 0 {
			z = Inf()
			return z
		}
		z = NegInf()
		return z
	}
	if b.cmp(NegInf()) == 0 || b.cmp(Inf()) == 0 {
		z = NewFrac(0, 1)
		return z
	}

	z.num.Mul(a.num, b.denom)
	z.denom.Mul(a.denom, b.num)
	return z.reduce()
}

//cmp compares a and b
//if a < b a.cmp(b) == -1
//if a > b a.cmp(b) == 1
//if a = b a.cmp(b) == 0
//Specific comparasions:
//
//WARINING: NaN().cmp(x) or x.cmp(NaN()) will lead to panic for every x except of NaN
//
//Inf().cmp(Inf()) == 0
//Inf().cmp(x) == 1 for any x
//x.cmp(Inf()) = -1 for any x
//NegInf().cmp(NegInf()) == 0 for any x
//NegInf().cmp(x) == -1 for any x
//x.cmp(NegInf()) = 1 for any x
//NaN().cmp(NaN) == 0 for any x. Read warning above.
func (a Value) cmp(b *Value) int {
	a.checkNil()
	b.checkNil()

	if a.denom.Sign() == 0 && b.denom.Sign() == 0 &&
		a.num.Sign() == 0 && b.num.Sign() == 0 {
		return 0
	}

	if a.denom.Sign() == 0 {
		if a.num.Sign() == 0 ||
			b.denom.Sign() == 0 && b.num.Sign() == 0 {
			panic("Comparing with NaN")
		}
		if a.num.Sign() < 0 {
			if b.denom.Sign() == 0 && b.num.Sign() < 0 {
				return 0
			}
			return -1
		}
		if b.denom.Sign() == 0 && b.num.Sign() > 0 {
			return 0
		}
		return 1
	}

	if b.denom.Sign() == 0 {
		if b.num.Sign() == 0 ||
			a.denom.Sign() == 0 && a.num.Sign() == 0 {
			panic("Comparing with NaN")
		}
		if b.num.Sign() < 0 {
			if a.denom.Sign() == 0 && a.num.Sign() < 0 {
				return 0
			}
			return 1
		}
		if a.denom.Sign() == 0 && a.num.Sign() > 0 {
			return 0
		}
		return -1
	}

	if a.sign() != b.sign() {
		if a.sign() < b.sign() {
			return -1
		}
		return 1
	}
	if a.sign() == 0 && b.sign() == 0 {
		return 0
	}
	tmp1 := new(big.Int)
	tmp2 := new(big.Int)
	if diff := new(big.Int).Sub(tmp1.Div(a.num, a.denom), tmp2.Div(b.num, b.denom)); diff.Sign() != 0 {
		return diff.Sign()
	} else {
		nod := nod(a.denom, b.denom)
		diff := new(big.Int).Sub(
			tmp1.Mul(a.num, new(big.Int).Div(b.denom, nod)),
			tmp2.Mul(b.num, new(big.Int).Div(a.denom, nod)),
		)
		return diff.Sign()
	}
}

func (v *Value) checkNil() {
	if v.num == nil {
		v.num = new(big.Int)
	}
	if v.denom == nil {
		v.denom = new(big.Int)
		v.denom.SetInt64(1)
	}
}

func nod(a *big.Int, b *big.Int) *big.Int {
	r := big.NewInt(1)
	prevR := new(big.Int).Set(r)
	for r.Sign() != 0 {
		prevR.Set(r)
		r.Mod(a, b)
		a = b
		b = new(big.Int).Set(r)
	}
	return prevR
}
