package domain

type interval interface {
	Operation
	Left() Value
	Right() Value
	Add(b interval) interval
	Sub(b interval) interval
	Div(b interval) interval
	Mul(b interval) interval
	String() string
}

type intervalImpl struct {
	a Value
	b Value
}

func (i intervalImpl) Left() Value {
	return i.a
}

func (i intervalImpl) Right() Value {
	return i.b
}

func (i intervalImpl) Add(b interval) interval {
	newA := i.a.Add(b.Left())
	newB := i.b.Add(b.Right())
	return intervalImpl{
		a: newA,
		b: newB,
	}
}

func (i intervalImpl) Sub(b interval) interval {
	newA := i.a.Sub(b.Left())
	newB := i.b.Sub(b.Right())
	return intervalImpl{
		a: newA,
		b: newB,
	}
}
func (i intervalImpl) Div(b interval) interval {
	newA := i.a.Div(b.Left())
	newB := i.b.Div(b.Right())
	return intervalImpl{
		a: newA,
		b: newB,
	}
}
func (i intervalImpl) Mul(b interval) interval {
	newA := i.a.Mul(b.Left())
	newB := i.b.Mul(b.Right())
	return intervalImpl{
		a: newA,
		b: newB,
	}
}

func (i intervalImpl) String() string {
	return "[ " + i.a.String() + ", " + i.b.String() + " ]"
}

func (i intervalImpl) Result() interval {
	return i
}

func Interval(left Value, right Value) interval {
	return intervalImpl{
		a: left,
		b: right,
	}
}
