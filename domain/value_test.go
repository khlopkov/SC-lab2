package domain

import "testing"

func TestValueCmp(t *testing.T) {
	sign := func(a int) int {
		if a == 0 {
			return 0
		}
		if a < 0 {
			return -1
		}
		return 1
	}
	var testPairs = []struct {
		a   *Value
		b   *Value
		res int
	}{
		{
			a:   NewFrac(1225, 210),
			b:   NewFrac(35, 6),
			res: 0,
		}, {
			a:   NewFrac(-1225, 210),
			b:   NewFrac(35, 6),
			res: -1,
		}, {
			a:   NewFrac(35, 6),
			b:   NewFrac(-1225, 210),
			res: 1,
		}, {
			a:   NewFrac(-1225, -210),
			b:   NewFrac(35, 6),
			res: 0,
		}, {
			a:   NewFrac(1225, -210),
			b:   NewFrac(35, 6),
			res: -1,
		}, {
			a:   NewFrac(1227, 210),
			b:   NewFrac(35, 6),
			res: 1,
		}, {
			a:   NewFrac(35, 6),
			b:   NewFrac(1227, 210),
			res: -1,
		}, {
			a:   NewFrac(1225, 212),
			b:   NewFrac(35, 6),
			res: -1,
		}, {
			a:   NewFrac(35, 6),
			b:   NewFrac(1225, 212),
			res: 1,
		}, {
			a:   NewFrac(2, 3),
			b:   NewFrac(15, 6),
			res: -1,
		}, {
			a:   NewFrac(15, 6),
			b:   NewFrac(2, 3),
			res: 1,
		}, {
			a:   NewFrac(2, 3),
			b:   NewFrac(1, 2),
			res: 1,
		}, {
			a:   NewFrac(1, 2),
			b:   NewFrac(2, 3),
			res: -1,
		},
		{
			a:   NewFrac(1, 2),
			b:   NewFrac(1, 0),
			res: -1,
		},
		{
			a:   NewFrac(1, 0),
			b:   NewFrac(1, 2),
			res: 1,
		},
		{
			a:   NewFrac(-1, 2),
			b:   NewFrac(1, 0),
			res: -1,
		},
		{
			a:   NewFrac(1, 0),
			b:   NewFrac(-1, 2),
			res: 1,
		},
		{
			a:   NewFrac(-1, 0),
			b:   NewFrac(1, 2),
			res: -1,
		},
		{
			a:   NewFrac(1, 2),
			b:   NewFrac(-1, 0),
			res: 1,
		},
		{
			a:   NewFrac(1, 0),
			b:   NewFrac(1, 0),
			res: 0,
		},
		{
			a:   NewFrac(0, 0),
			b:   NewFrac(0, 0),
			res: 0,
		},
		{
			a:   NewFrac(-1, 0),
			b:   NewFrac(1, 0),
			res: -1,
		},
		{
			a:   NewFrac(1, 0),
			b:   NewFrac(-1, 0),
			res: 1,
		},
		{
			a:   NewFrac(0, 1),
			b:   NewFrac(0, -1),
			res: 0,
		},
	}
	for i, pair := range testPairs {
		if sign(pair.a.cmp(pair.b)) != pair.res {
			if pair.res == 0 {
				t.Errorf("In pair %d: %s should be equal %s", i, pair.a, pair.b)
			}
			if pair.res < 0 {
				t.Errorf("In pair %d: %s should be less then %s", i, pair.a, pair.b)
			}
			if pair.res > 0 {
				t.Errorf("In pair %d: %s should be more then %s", i, pair.a, pair.b)
			}
		}
	}
}

func TestStringValue(t *testing.T) {
	var testPairs = []struct {
		a   *Value
		res string
	}{
		{
			a:   NewFrac(35, 6),
			res: "35 / 6",
		},
		{
			a:   NewFrac(-35, 6),
			res: "-35 / 6",
		},
		{
			a:   NewFrac(35, -6),
			res: "-35 / 6",
		},
		{
			a:   NewFrac(-35, -6),
			res: "35 / 6",
		},
		{
			a:   NewFrac(0, -6),
			res: "0",
		},
		{
			a:   NewFrac(1, 0),
			res: "Inf",
		},
		{
			a:   NewFrac(0, 0),
			res: "NaN",
		},
		{
			a:   NewFrac(-1, 0),
			res: "-Inf",
		},
	}
	for i, pair := range testPairs {
		if pair.a.String() != pair.res {
			t.Errorf("In pair %d: %s should be equal %s", i, pair.a, pair.res)
		}
	}
}

func TestValueOperations(t *testing.T) {
	var testPairs = []struct {
		operation *Value
		res       *Value
	}{
		{
			operation: new(Value).add(NewFrac(12, 5), NewFrac(13, 5)),
			res:       NewFrac(5, 1),
		},
		{
			operation: new(Value).add(NewFrac(12, 7), NewFrac(13, 9)),
			res:       NewFrac(199, 63),
		},
		{
			operation: new(Value).add(NewFrac(12, 7), NewFrac(-13, 9)),
			res:       new(Value).sub(NewFrac(12, 7), NewFrac(13, 9)),
		},
		{
			operation: new(Value).add(NewFrac(12, 7), NewFrac(-13, 9)),
			res:       NewFrac(17, 63),
		},
		{
			operation: new(Value).add(NewFrac(12, 0), NewFrac(-13, 9)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(12, 5), NewFrac(13, 5)),
			res:       NewFrac(156, 25),
		},
		{
			operation: new(Value).mul(NewFrac(12, 7), NewFrac(13, 9)),
			res:       NewFrac(52, 21),
		},
		{
			operation: new(Value).mul(NewFrac(12, 7), NewFrac(-13, 9)),
			res:       NewFrac(-52, 21),
		},
		{
			operation: new(Value).mul(NewFrac(12, 7), NewFrac(9, 13)),
			res:       NewFrac(108, 91),
		},
		{
			operation: new(Value).mul(NewFrac(12, 7), NewFrac(9, 13)),
			res:       new(Value).div(NewFrac(12, 7), NewFrac(13, 9)),
		},
		{
			operation: new(Value).add(NewFrac(-12, 1), NewFrac(-13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).add(NewFrac(12, 1), NewFrac(-13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).add(NewFrac(-12, 0), NewFrac(-13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).add(NewFrac(12, 1), NewFrac(13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).add(NewFrac(12, 0), NewFrac(13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).add(NewFrac(0, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).add(NewFrac(1, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).add(NewFrac(1, 1), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).add(NewFrac(-12, 0), NewFrac(13, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).mul(NewFrac(-12, 1), NewFrac(-13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(12, 1), NewFrac(-13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(-12, 0), NewFrac(-13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(12, 1), NewFrac(13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(12, 0), NewFrac(13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).mul(NewFrac(0, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).mul(NewFrac(1, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).mul(NewFrac(1, 1), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).mul(NewFrac(-12, 0), NewFrac(13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).sub(NewFrac(12, 1), NewFrac(-13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).sub(NewFrac(-12, 0), NewFrac(-13, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).sub(NewFrac(12, 1), NewFrac(13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).sub(NewFrac(12, 0), NewFrac(13, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).sub(NewFrac(0, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).sub(NewFrac(1, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).sub(NewFrac(1, 1), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).sub(NewFrac(-12, 0), NewFrac(13, 0)),
			res:       NewFrac(-1, 0),
		},
		{
			operation: new(Value).sub(NewFrac(-12, 1), NewFrac(-13, 0)),
			res:       NewFrac(1, 0),
		},
		{
			operation: new(Value).div(NewFrac(12, 1), NewFrac(-13, 0)),
			res:       NewFrac(0, 1),
		},
		{
			operation: new(Value).div(NewFrac(-12, 1), NewFrac(-13, 0)),
			res:       NewFrac(0, 1),
		},
		{
			operation: new(Value).div(NewFrac(-12, 0), NewFrac(-13, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).div(NewFrac(12, 1), NewFrac(13, 0)),
			res:       NewFrac(0, 1),
		},
		{
			operation: new(Value).div(NewFrac(12, 0), NewFrac(13, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).div(NewFrac(0, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).div(NewFrac(1, 0), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).div(NewFrac(1, 1), NewFrac(0, 0)),
			res:       NewFrac(0, 0),
		},
		{
			operation: new(Value).div(NewFrac(-12, 0), NewFrac(13, 0)),
			res:       NewFrac(0, 0),
		},
	}
	for i, pair := range testPairs {
		if pair.operation.cmp(pair.res) != 0 {
			t.Errorf("In pair %d: %s should be equal %s", i, pair.operation, pair.res)
		}
	}
}
