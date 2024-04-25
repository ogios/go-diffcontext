package diffcontext

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
	// AdditionalSublines [2]bool // wether `Before` or `After` have multiple lines. 0-1, 1 means two lines
}

type DiffConstractor struct {
	Lines []*DiffLine
}

func New() *DiffConstractor {
	return &DiffConstractor{}
}

func (d *DiffConstractor) AddDiffs(ds []diffmatchpatch.Diff) {
	c := newConstractor()

	// for every Diff
	for i, d2 := range ds {
		// split lines
		lines := strings.Split(d2.Text, "\n")
		// for every line
		for i2, l := range lines {
			// if i2 == 0 && len(lines) > 1 && len(q) > 0 {
			if i2 == 0 && len(lines) > 1 && c.length > 0 {
				// the first one && lines length not 1 -- (means the end of one line)

				// add to queue and "try" resolve this line (markQ)
				c.addQ(diffData{
					diffState: d2.Type,
					data:      []byte(l),
				}, 0)
				c.markQ(d2.Type)
			} else if i2 == len(lines)-1 {
				// the end of lines(includes line length=1) -- (means) start of the new line but not end
				if l == "" && i != len(ds)-1 {
					// normally if lines length > 1, the end of lines may be just empty string("") meaning it's the start of new Line
					// but it should not be recorded in to queue which will result in wrong state(please check func `resolveQueue` where `DiffChanged` is computed), so just ignore it
					// but not the end of Diffs which would cause loosing one line at the end of content
					continue
				}
				// think of this: a Diff with DiffDelete + "anything\n\tsomething", this splits into 2 parts: `anything` & `\tsomething`
				// part 1 jumped into `if` judgment and run markQ which set state for c.qs[0] `delete`(4)
				// and now part 2 jumped into here, we have to make another queue(c.qs[1]) to temprarily save this since c.qs[0] is not finished yet
				//
				// now here comes another Diff with DiffInsert + "anything2\n\tsomething2", this also splits into 2 parts: `anything2` & `\tsomething2`
				// part 1 jumped into `if` judgment and run markQ which set state to 10(4+6), resolve c.qs[0] and move q.cs[1] to q.cs[0]
				// and part 2 jump into here, now c.qs[0] is not in any state, we need to push part 2 into c.qs[0], not creating another one
				var index int
				if (len(lines) > 1 && c.state > 0) || c.length == 0 {
					index = c.makeNewQ()
				}
				c.addQ(diffData{
					diffState: d2.Type,
					data:      []byte(l),
				}, index)

			} else {
				// normal Lines
				// basically every thing about queue is on the first and last line
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
				c.addLine(&DiffLine{
					Before: b,
					After:  a,
					State:  d2.Type,
				})
			}
		}
		if i == len(ds)-1 {
			// the end of content and queue is not empty
			for c.length > 0 {
				c.resQ()
			}
		}
	}
	d.Lines = append(d.Lines, c.dLines...)
}

const (
	LINE_BREAK = byte('\n')
	EQUAL      = byte(' ')
	DEL        = byte('-')
	INS        = byte('+')
)

func getBefore(lines []*DiffLine, withFront bool) []byte {
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

func getAfter(lines []*DiffLine, withFront bool) []byte {
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
