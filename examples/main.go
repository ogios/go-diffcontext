package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/andreyvit/diff"
	"github.com/ogios/go-diffcontext"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	test()
	benchMark()
}

func test() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(code1, code2, true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupSemantic(diffs)

	dc := diffcontext.New()
	dc.AddDiffs(diffs)
	fmt.Printf("dc.Lines: %v\n", dc.Lines)
	fmt.Printf("dc.GetBefore(): %v\n", dc.GetBefore())
	fmt.Printf("dc.GetAfter(): %v\n", dc.GetAfter())
	fmt.Printf("dc.GetMixed(): %v\n", dc.GetMixed())
	fmt.Printf("diff.LineDiff(code1, code2): %v\n", diff.LineDiff(code1, code2))
}

func benchMark() {
	marks := 1000
	c := make([]string, marks)
	runtime.GC()
	start := time.Now()
	for i := 0; i < marks; i++ {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(code1, code2, true)
		diffs = dmp.DiffCleanupSemantic(diffs)
		diffs = dmp.DiffCleanupSemantic(diffs)
		dc := diffcontext.New()
		dc.AddDiffs(diffs)
		c[i] = dc.GetMixed()
	}
	fmt.Println(time.Now().UnixMilli() - start.UnixMilli())

	c = make([]string, marks)
	runtime.GC()
	start = time.Now()
	for i := 0; i < marks; i++ {
		c[i] = diff.LineDiff(code1, code2)
	}
	fmt.Println(time.Now().UnixMilli() - start.UnixMilli())
}
