package main

import (
	"fmt"
	"go-diff-test/diffconstract"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(code1, code2, true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupSemantic(diffs)

	dc := diffconstract.New()
	dc.AddDiffs(diffs)
	fmt.Printf("dc.Lines: %v\n", dc.Lines)
	fmt.Printf("dc.GetBefore(): %v\n", dc.GetBefore())
	fmt.Printf("dc.GetAfter(): %v\n", dc.GetAfter())

	// fmt.Println(dmp.DiffPrettyText(diffs))
	// fmt.Println(diffs)
}
