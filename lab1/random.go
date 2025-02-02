package main

import "math"

const (
	a = 5
	b = 8
	c = 1
)

var (
	m = int(math.Pow(2, b))
)

type PseudoRandom struct {
	t int
}

func NewRandom() *PseudoRandom {
	return &PseudoRandom{t: 2}
}

func (r *PseudoRandom) Next() int {
	r.t = (a*r.t + c) % m
	return r.t
}
