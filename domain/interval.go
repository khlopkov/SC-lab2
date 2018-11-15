package domain

import (
	"errors"
	"regexp"
	"strconv"
)

var varRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

type AbstractInterval interface {
	Operation
	Left() (float64, bool)
	Right() (float64, bool)
}

type intervalImpl struct {
	left  float64
	right float64
}

func (i intervalImpl) Left() (float64, bool) {
	return i.left, true
}

func (i intervalImpl) Right() (float64, bool) {
	return i.right, true
}

func (i intervalImpl) String() string {
	return "[ " + strconv.FormatFloat(i.left, 'f', 4, 64) + ", " + strconv.FormatFloat(i.right, 'f', 4, 64) + " ]"
}

func (i intervalImpl) Solve(params ParamMap) Operation {
	return i
}

func (i intervalImpl) priority() byte {
	return 255
}

func Interval(left float64, right float64) AbstractInterval {
	return intervalImpl{
		left:  left,
		right: right,
	}
}

type parametricInterval struct {
	variable string
}

func (i parametricInterval) Left() (float64, bool) {
	return 0, false
}

func (i parametricInterval) Right() (float64, bool) {
	return 0, false
}

func (i parametricInterval) String() string {
	return i.variable
}

func (i parametricInterval) Solve(params ParamMap) Operation {
	if params[i.variable] != nil {
		return params[i.variable]
	} else {
		return i
	}
}

func (i parametricInterval) priority() byte {
	return 255
}

func ParametricInterval(varName string) (AbstractInterval, error) {
	if !varRegexp.MatchString(varName) {
		return nil, errors.New("wrong variable name. Variable should contain letters and numbers, starting from letter")
	}
	return parametricInterval{
		variable: varName,
	}, nil
}
