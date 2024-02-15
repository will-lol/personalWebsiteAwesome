package eid

import (
	"strconv"
)

type EidFactory struct {
	state int64
}

func NewEidFactory() EidFactory {
	return EidFactory{state: 0}
}

func (e EidFactory) GetNext() string {
	eid := strconv.FormatInt(e.state, 36)
	e.state += 1

	return eid
}

