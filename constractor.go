package diffcontext

import (
	"errors"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type diffData struct {
	data      []byte
	diffState diffmatchpatch.Operation
}

type constractorQ struct {
	q   []diffData
	pos int
}

type constractor struct {
	qs     [2]*constractorQ // mianly operating on qs[0], qs[1] is just for temprary saving
	dLines []DiffLine       // lines
	length int              // refers to qs
	state  int
}

func newConstractor() constractor {
	return constractor{
		dLines: make([]DiffLine, 0),
	}
}

func (c *constractor) makeNewQ() int {
	if c.length >= 2 {
		panic(errors.New("max 2 queue"))
	}
	c.dLines = append(c.dLines, DiffLine{})
	con := &constractorQ{
		q:   make([]diffData, 0),
		pos: len(c.dLines) - 1,
	}
	if c.qs[0] != nil {
		c.qs[1] = con
	} else {
		c.qs[0] = con
	}
	c.length++
	return c.length - 1
}

func (c *constractor) resQ() {
	q := c.qs[0]
	c.dLines[q.pos] = resolveQueue(q.q)
	c.qs[0] = nil
	if c.qs[1] != nil {
		c.qs[0] = c.qs[1]
		c.qs[1] = nil
	}
	c.length--
}

func (c *constractor) addQ(d diffData, i int) {
	c.qs[i].q = append(c.qs[i].q, d)
}

func (c *constractor) markQ(t diffmatchpatch.Operation) {
	switch t {
	case diffmatchpatch.DiffEqual:
		c.state += 10
	case diffmatchpatch.DiffInsert:
		if c.state < 6 {
			c.state += 6
		} else {
			panic(fmt.Errorf("state error: %d %s", c.state, "insert"))
		}
	case diffmatchpatch.DiffDelete:
		if c.state < 4 {
			c.state += 4
		} else {
			panic(fmt.Errorf("state error: %d %s", c.state, "delete"))
		}
	}
	if c.state >= 10 {
		c.resQ()
		c.state = 0
	}
}

func (c *constractor) addLine(l DiffLine) {
	c.dLines = append(c.dLines, l)
}
