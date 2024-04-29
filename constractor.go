package diffcontext

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

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
	dLines []*DiffLine      // lines
	length int              // refers to qs
	state  int
}

func newConstractor() constractor {
	return constractor{
		dLines: make([]*DiffLine, 0),
	}
}

func (c *constractor) makeNewQ() int {
	if c.length >= 2 {
		panic(errors.New("max 2 queue"))
	}
	c.dLines = append(c.dLines, &DiffLine{})
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

var line_break_bytes = []byte{line_break}

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
	if c.state == 10 {
		c.resQ()
		c.state = 0
	} else if c.state > 10 {
		if c.qs[1] != nil {
			// newQ := make([]diffData, len(c.qs[0].q)+len(c.qs[1].q)+1)
			// first := c.qs[0].q
			// second := c.qs[1].q
			// i := len(first) - 1
			// last := first[i]
			// copy(newQ, first[:i])
			// if c.state-10 == 6 {
			// 	newQ[i] = diffData{
			// 		data:      line_break_bytes,
			// 		diffState: diffmatchpatch.DiffInsert,
			// 	}
			// } else {
			// 	newQ[i] = diffData{
			// 		data:      line_break_bytes,
			// 		diffState: diffmatchpatch.DiffDelete,
			// 	}
			// }
			// i++
			// copy(newQ[i:], second)
			// newQ[len(newQ)-1] = last
			// c.qs[0].q = newQ
			// c.dLines = slices.Delete(c.dLines, c.qs[1].pos, c.qs[1].pos+1)
			// c.qs[1] = nil
			// c.length--
			// c.resQ()
			// c.state = 0
			first := c.qs[0].q
			second := c.qs[1].q
			if c.state-10 == 6 {
				first[len(first)-1].diffState = diffmatchpatch.DiffDelete
			} else {
				first[len(first)-1].diffState = diffmatchpatch.DiffInsert
			}
			c.qs[1].q = append(second, diffData{
				data:      first[len(first)-1].data,
				diffState: second[0].diffState,
			})
			c.resQ()
			c.resQ()
			c.state = 0
		} else {
			first := c.qs[0].q
			last := first[len(first)-1]
			if c.state-10 == 6 {
				last.diffState = diffmatchpatch.DiffInsert
			} else {
				last.diffState = diffmatchpatch.DiffDelete
			}
			c.addQ(last, c.makeNewQ())
			c.qs[0].q = slices.Delete(first, len(first)-1, len(first))
			c.resQ()
			c.resQ()
			c.state = 0
		}
	}
}

func (c *constractor) addLine(l *DiffLine) {
	c.dLines = append(c.dLines, l)
}

func resolveQueue(q []diffData) *DiffLine {
	var before, after bytes.Buffer
	state := q[0].diffState
	setState := func(s diffmatchpatch.Operation) {
		if state != s {
			state = DiffChanged
		}
	}
	// var lineBreakState diffmatchpatch.Operation
	for _, dd := range q {
		// if len(dd.data) > 0 {
		// 	if dd.data[0] == LINE_BREAK {
		// 		lineBreakState = dd.diffState
		// 	}
		// }
		switch dd.diffState {
		case diffmatchpatch.DiffEqual:
			setState(diffmatchpatch.DiffEqual)
			before.Write(dd.data)
			after.Write(dd.data)
		case diffmatchpatch.DiffDelete:
			setState(diffmatchpatch.DiffDelete)
			before.Write(dd.data)
		case diffmatchpatch.DiffInsert:
			setState(diffmatchpatch.DiffInsert)
			after.Write(dd.data)
		}
	}
	data := &DiffLine{
		Before: before.Bytes(),
		After:  after.Bytes(),
		State:  state,
	}
	// if lineBreakState == diffmatchpatch.DiffDelete {
	// 	data.AdditionalSublines[0] = true
	// } else {
	// 	data.AdditionalSublines[1] = true
	// }
	return data
}
