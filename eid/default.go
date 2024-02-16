package eid

import (
	"strconv"
	"sync"
)

type EidFactory struct {
	state int64
	m sync.Mutex
}

func NewEidFactory() EidFactory {
	return EidFactory{state: 0}
}

func (e *EidFactory) GetNext() string {
	e.m.Lock()
	defer e.m.Unlock()
	eid := strconv.FormatInt(e.state, 36)
	e.state = e.state + 1

	return eid
}

