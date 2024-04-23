package diffconstract

import (
	"bytes"
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
	Lines []DiffLine
	q     []diffData
}

func New() *DiffConstractor {
	return &DiffConstractor{
		q: make([]diffData, 0),
	}
}

func (d *DiffConstractor) AddDiffs(ds []diffmatchpatch.Diff) {
	for i, d2 := range ds {
		lines := strings.Split(d2.Text, "\n")
		for i2, l := range lines {
			if i2 == 0 && len(lines) > 1 && len(d.q) > 0 {
				d.q = append(d.q, diffData{
					diffState: d2.Type,
					data:      []byte(l),
				})
				d.resolveQueue()
			} else if i2 == len(lines)-1 {
				if l == "" && i != len(ds)-1 {
					continue
				}
				d.q = append(d.q, diffData{
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
				data := DiffLine{
					Before: b,
					After:  a,
					State:  d2.Type,
				}
				d.Lines = append(d.Lines, data)
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
	data := DiffLine{
		Before: before.Bytes(),
		After:  after.Bytes(),
		State:  state,
	}
	d.Lines = append(d.Lines, data)
	d.q = make([]diffData, 0)
}

const (
	LINE_BREAK = byte('\n')
	EQUAL      = byte(' ')
	DEL        = byte('-')
	INS        = byte('+')
)

func getBefore(lines []DiffLine, withFront bool) []byte {
	var builder bytes.Buffer
	// var builder strings.Builder
	for i, dl := range lines {
		be := dl.Before
		if dl.State != diffmatchpatch.DiffInsert {
			if withFront {
				if dl.State != diffmatchpatch.DiffEqual {
					builder.WriteByte(DEL)
				} else {
					builder.WriteByte(EQUAL)
				}
			}
			builder.Write(be)
			if i != len(lines)-1 {
				builder.WriteByte(LINE_BREAK)
			}
		}
	}
	return builder.Bytes()
}

func getAfter(lines []DiffLine, withFront bool) []byte {
	var builder bytes.Buffer
	// var builder strings.Builder
	for i, dl := range lines {
		af := dl.After
		if dl.State != diffmatchpatch.DiffDelete {
			if withFront {
				if dl.State != diffmatchpatch.DiffEqual {
					builder.WriteByte(INS)
				} else {
					builder.WriteByte(EQUAL)
				}
			}
			builder.Write(af)
			if i != len(lines)-1 {
				builder.WriteByte(LINE_BREAK)
			}
		}
	}
	return builder.Bytes()
}

func (d *DiffConstractor) GetBefore() string {
	return string(getBefore(d.Lines, false))
}

func (d *DiffConstractor) GetAfter() string {
	return string(getAfter(d.Lines, false))
}

func (d *DiffConstractor) GetMixed() string {
	var builder strings.Builder
	inChange := -1
	for i, dl := range d.Lines {
		if inChange < 0 {
			if dl.State == diffmatchpatch.DiffEqual {
				builder.WriteByte(EQUAL)
				builder.Write(dl.Before)
				if i != len(d.Lines)-1 {
					builder.WriteByte(LINE_BREAK)
				}
			} else {
				inChange = i
			}
		} else {
			if dl.State == diffmatchpatch.DiffEqual {
				changes := d.Lines[inChange:i]
				before := getBefore(changes, true)
				if len(before) > 0 {
					builder.Write(before)
					builder.WriteByte(LINE_BREAK)
				}
				after := getAfter(changes, true)
				if len(after) > 0 {
					builder.Write(after)
					builder.WriteByte(LINE_BREAK)
				}
				builder.WriteByte(EQUAL)
				builder.Write(dl.Before)
				if i != len(d.Lines)-1 {
					builder.WriteByte(LINE_BREAK)
				}
				inChange = -1
			}
		}
	}
	if inChange >= 0 {
		builder.WriteByte(LINE_BREAK)
		changes := d.Lines[inChange:]
		builder.Write(getBefore(changes, true))
		builder.WriteByte(LINE_BREAK)
		builder.Write(getAfter(changes, true))
	}
	return builder.String()
}
