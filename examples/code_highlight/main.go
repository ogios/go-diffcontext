package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/andreyvit/diff"
	"github.com/ogios/go-diffcontext"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	dc := test()
	matchLine(dc)
}

func test() *diffcontext.DiffConstractor {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(code1, code2, true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupEfficiency(diffs)

	dc := diffcontext.New()
	dc.AddDiffs(diffs)
	fmt.Printf("dc.Lines:\n%v\n", dc.Lines)
	fmt.Printf("dc.GetBefore():\n%v\n", dc.GetBefore())
	fmt.Printf("dc.GetAfter():\n%v\n", dc.GetAfter())
	fmt.Printf("dc.GetMixed():\n%v\n", dc.GetMixed())
	fmt.Printf("diff.LineDiff(code1, code2):\n%v\n", diff.LineDiff(code1, code2))
	return dc
}

func matchLine(dc *diffcontext.DiffConstractor) {
	c1 := highlight(code1)
	linesC1 := strings.Split(c1, "\n")
	// linesC1 := strings.Split(code1, "\n")

	c2 := highlight(code2)
	linesC2 := strings.Split(c2, "\n")
	// linesC2 := strings.Split(code2, "\n")
	i1 := 0
	i2 := 0
	for _, dl := range dc.Lines {
		be := []byte(linesC1[i1])
		af := []byte(linesC2[i2])
		switch dl.State {
		case diffmatchpatch.DiffEqual:
			dl.Before, dl.After = be, be
			// fmt.Println(string(dl.After) == string(dl.Before), string(dl.After) == string(be))
			i1++
			i2++
		case diffcontext.DiffChanged:
			dl.Before, dl.After = be, af
			// fmt.Println(string(dl.After) == string(af), string(dl.Before) == string(be))
			i1++
			i2++
		case diffmatchpatch.DiffInsert:
			dl.After = af
			// fmt.Println(string(dl.After) == string(af))
			i2++
		case diffmatchpatch.DiffDelete:
			dl.Before = be
			// fmt.Println(string(dl.Before) == string(be))
			i1++
		}
	}
	fmt.Println(dc.GetMixed())
}

func highlight(src string) string {
	buf := new(bytes.Buffer)
	err := quick.Highlight(buf, src, "go", "terminal16m", "catppuccin-mocha")
	if err != nil {
		panic(err)
	}
	return buf.String()
}
