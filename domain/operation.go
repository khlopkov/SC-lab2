package domain

type Operation interface {
	Result() interval
}

type BinaryOperation interface {
	Operation
	A() Operation
	B() Operation
	SetA(o Operation)
	SetB(o Operation)
}

type AbstractBinaryOperation struct {
	a Operation
	b Operation
}

func (s AbstractBinaryOperation) A() Operation {
	return s.a
}
func (s AbstractBinaryOperation) B() Operation {
	return s.b
}
func (s *AbstractBinaryOperation) SetA(o Operation) {
	s.a = o
}
func (s *AbstractBinaryOperation) SetB(o Operation) {
	s.b = o
}

type add struct {
	AbstractBinaryOperation
}

func (s add) Result() interval {
	return s.b.Result().Add(s.a.Result())
}

func Add(a Operation, b Operation) Operation {
	return add{
		AbstractBinaryOperation{
			a: a,
			b: b,
		},
	}
}

type sub struct {
	AbstractBinaryOperation
}

func (s sub) Result() interval {
	return s.b.Result().Sub(s.a.Result())
}

func Sub(a Operation, b Operation) Operation {
	return sub{
		AbstractBinaryOperation{
			a: a,
			b: b,
		},
	}
}

type mul struct {
	AbstractBinaryOperation
}

func (s mul) Result() interval {
	return s.b.Result().Mul(s.a.Result())
}

func Mul(a Operation, b Operation) Operation {
	return mul{
		AbstractBinaryOperation{
			a: a,
			b: b,
		},
	}
}

type div struct {
	AbstractBinaryOperation
}

func (s div) Result() interval {
	return s.b.Result().Div(s.a.Result())
}

func Div(a Operation, b Operation) Operation {
	return div{
		AbstractBinaryOperation{
			a: a,
			b: b,
		},
	}
}
