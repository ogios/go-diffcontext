package main

import (
	"fmt"
	"os"

	"github.com/andreyvit/diff"
	"github.com/ogios/go-diffcontext"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	test()
}

func test() {
	dmp := diffmatchpatch.New()
	code1, _ := os.ReadFile("./code1.txt")
	code2, _ := os.ReadFile("./code2.txt")
	// code1, _ := os.ReadFile("./test1")
	// code2, _ := os.ReadFile("./test2")
	fmt.Printf("diff.LineDiff(code1, code2): %v\n", diff.LineDiff(string(code1), string(code2)))

	diffs := dmp.DiffMain(string(code1), string(code2), true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupEfficiency(diffs)

	dc := diffcontext.New()
	dc.AddDiffs(diffs)
	fmt.Printf("dc.Lines: %v\n", dc.Lines)
	fmt.Printf("dc.GetBefore(): %v\n", dc.GetBefore())
	fmt.Printf("dc.GetAfter(): %v\n", dc.GetAfter())
	fmt.Printf("dc.GetMixed(): %v\n", dc.GetMixed())
}
