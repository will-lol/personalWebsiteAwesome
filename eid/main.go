package eid

import (
	"strconv"
)

type EidFactory struct {
	state int
}

func New() EidFactory {
	return EidFactory{state: 0}
}

func (e EidFactory) CreateEid() string {
	eid := strconv.Itoa(e.state)
	e.state += 1
	return eid 
}
