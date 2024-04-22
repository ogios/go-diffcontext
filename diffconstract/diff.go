package diffconstract

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const DiffChanged diffmatchpatch.Operation = 2

type DiffLine struct {
	Before []byte
	After  []byte
	State  diffmatchpatch.Operation
}

type diffData struct {
	data      []byte
	diffState diffmatchpatch.Operation
}

type DiffConstractor struct {
	Lines []*DiffLine
	q     []*diffData
}

func New() *DiffConstractor {
	return &DiffConstractor{
		q: make([]*diffData, 0),
	}
}

func endWithBreak(lines []string) bool {
	l := len(lines)
	return l > 1 && lines[l-1] == ""
}

func (d *DiffConstractor) AddDiffs(ds []diffmatchpatch.Diff) {
	for i, d2 := range ds {
		lines := strings.Split(d2.Text, "\n")
		for i2, l := range lines {
			if i2 == 0 && len(lines) > 1 && len(d.q) > 0 {
				d.q = append(d.q, &diffData{
					diffState: d2.Type,
					data:      []byte(l),
				})
				d.resolveQueue()
			} else if i2 == len(lines)-1 && !endWithBreak(lines) {
				d.q = append(d.q, &diffData{
					diffState: d2.Type,
					data:      []byte(l),
				})
			} else {
				var b, a []byte
				switch d2.Type {
				case diffmatchpatch.DiffEqual:
					temp := []byte(l)
					b = temp
					a = temp
				case diffmatchpatch.DiffDelete:
					b = []byte(l)
				case diffmatchpatch.DiffInsert:
					a = []byte(l)
				}
				data := &DiffLine{
					Before: b,
					After:  a,
					State:  d2.Type,
				}
				d.Lines = append(d.Lines, data)
				fmt.Println("add line:", data)

			}
		}
		if i == len(ds)-1 && len(d.q) != 0 {
			d.resolveQueue()
		}
	}
}

func (d *DiffConstractor) resolveQueue() {
	var before, after bytes.Buffer
	state := d.q[0].diffState
	setState := func(s diffmatchpatch.Operation) {
		if state != s {
			state = DiffChanged
		}
	}
	for _, dd := range d.q {
		switch dd.diffState {
		case diffmatchpatch.DiffEqual:
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
	d.Lines = append(d.Lines, data)
	fmt.Println("resolve queue:", data)
	d.q = make([]*diffData, 0)
}

const LINE_BREAK = byte('\n')

func (d *DiffConstractor) GetBefore() string {
	var builder strings.Builder
	for i, dl := range d.Lines {
		be := dl.Before
		if dl.State != diffmatchpatch.DiffInsert {
			builder.Write(be)
			if i != len(d.Lines)-1 {
				builder.WriteByte(LINE_BREAK)
			}
		}
	}
	return builder.String()
}

func (d *DiffConstractor) GetAfter() string {
	var builder strings.Builder
	for i, dl := range d.Lines {
		af := dl.After
		if dl.State != diffmatchpatch.DiffDelete {
			builder.Write(af)
			if i != len(d.Lines)-1 {
				builder.WriteByte(LINE_BREAK)
			}
		}
	}
	return builder.String()
}
