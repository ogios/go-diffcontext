package diffcontext

import (
	"bytes"
	"errors"
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

func (c *constractor) resQ(i int) {
	q := c.qs[i]
	c.dLines[q.pos] = resolveQueue(q.q)
	c.qs[i] = nil
	if c.qs[1] != nil {
		c.qs[0] = c.qs[1]
		c.qs[1] = nil
	}
	c.length--
}

var line_break_bytes = []byte{line_break}

func (c *constractor) addQ(d diffData) {
	if c.qs[0] == nil {
		c.makeNewQ()
	}
	var to0, to1 diffmatchpatch.Operation
	switch c.state {
	case 6:
		to0 = diffmatchpatch.DiffDelete
		to1 = diffmatchpatch.DiffInsert
	case 4:
		to0 = diffmatchpatch.DiffInsert
		to1 = diffmatchpatch.DiffDelete
	case 0:
		c.qs[0].q = append(c.qs[0].q, d)
		return
	}
	switch d.diffState {
	case diffmatchpatch.DiffEqual:
		c.qs[0].q = append(c.qs[0].q, diffData{
			data:      d.data,
			diffState: to0,
		})
		c.qs[1].q = append(c.qs[1].q, diffData{
			data:      d.data,
			diffState: to1,
		})
	case to1:
		c.qs[1].q = append(c.qs[1].q, diffData{
			data:      d.data,
			diffState: to1,
		})
	case to0:
		c.qs[0].q = append(c.qs[0].q, diffData{
			data:      d.data,
			diffState: to0,
		})
	}
}

func (c *constractor) markQ(t diffmatchpatch.Operation) {
	resolveLast := func() {
		pos := c.qs[1].pos + 1
		c.resQ(1)
		c.qs[1] = &constractorQ{
			q:   make([]diffData, 0),
			pos: pos,
		}
		c.dLines = slices.Insert(c.dLines, pos, &DiffLine{})
		c.length++
	}
	switch t {
	case diffmatchpatch.DiffEqual:
		c.state += 10
	case diffmatchpatch.DiffInsert:
		if c.state == 0 || c.state == 4 {
			c.state += 6
			c.makeNewQ()
		} else {
			resolveLast()
		}
	case diffmatchpatch.DiffDelete:
		if c.state == 0 || c.state == 6 {
			c.state += 4
			c.makeNewQ()
		} else {
			resolveLast()
		}
	}

	if c.state >= 10 {
		for c.length > 0 {
			c.resQ(0)
		}
		c.state = 0
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
