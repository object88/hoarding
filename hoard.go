package hoarding

import (
	"fmt"
	"io"
	"sync"
)

type Level int

const (
	Always Level = iota
	Error
	Info
)

// type Hoarder interface {
// 	io.Closer

// 	Flush() error
// 	Msgf(l Level, format string, args ...interface{})
// }

type Hoard struct {
	cond       func() bool
	downstream io.Writer

	msgs  []message
	msgsL sync.Mutex
}

type message struct {
	l Level
	m string
}

func NewHoarder(downstream io.Writer, cond func() bool) *Hoard {
	return &Hoard{
		downstream: downstream,
	}
}

func (h *Hoard) Flush() error {
	if h == nil {
		return nil
	}

	h.msgsL.Lock()
	for _, m := range h.msgs {
		if m.l < m.l {
			continue
		}
		h.downstream.Write([]byte(m.m))
	}
	h.msgs = nil
	h.msgsL.Unlock()

	return nil
}

// Close satisfies the io.Closer interface
func (h *Hoard) Close() error {
	h.Flush()

	if h != nil {
		h.downstream = nil
	}

	return nil
}

func (h *Hoard) Msgf(l Level, format string, args ...interface{}) {
	if h == nil {
		return
	}

	m := format
	if len(args) > 0 {
		m = fmt.Sprintf(m, args...)
	}

	h.msgsL.Lock()
	h.msgs = append(h.msgs, message{
		l: l,
		m: m,
	})
	h.msgsL.Unlock()
}
