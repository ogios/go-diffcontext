# go-diffcontext

> Based on [`sergi/go-diff`](https://github.com/sergi/go-diff)

Transform from

```go
// Diff represents one diff operation
type Diff struct {
	Type Operation
	Text string
}
```

To

```go
type DiffLine struct {
	Before []byte
	After  []byte
	State  diffmatchpatch.Operation
}
```

## State
There were 3 so called `Operation`s in [`sergi/go-diff`](https://github.com/sergi/go-diff) which are:
```
const (
// DiffDelete item represents a delete diff.
DiffDelete Operation = -1
// DiffInsert item represents an insert diff.
DiffInsert Operation = 1
// DiffEqual item represents an equal diff.
DiffEqual Operation = 0
```
But in here we added one more state to clarify ***line change*** and ***slice of content change inside a line***
- if a line is removed or inserted, the `DiffLine` will be marked as `diffmatchpatch.DiffDelete` or `diffmatchpatch.DiffInsert` (package `diffmatchpatch` is from [`sergi/go-diff`](https://github.com/sergi/go-diff))
- if slice of string in one line changed, the `DiffLine` will be marked as `diffline.DiffChanged` (package `diffline` is from this repo), and both ***line content before change*** and ***line content after change*** are in `DiffLine.Before` and `DiffLine.After`

## Available funcs
Able to get both content before and after change by `GetBefore` and `GetAfter`

Also able to get mixed diff content with `GetMixed`

```
dc.GetMixed():
 package main

 import (
 	"fmt"

 	"github.com/sergi/go-diff/diffmatchpatch"
 )

 const (
 	text1 = "Lorem ipsum dolor."
-	text2 = "Lorem dolor sit amet."
+	text2 = "Lorem sit amet."
 )

 func main() {
-	dmp := diffmatchpatch.New()
-
 	diffs := dmp.DiffMain(text1, text2, false)

+	dmp := diffmatchpatch.New()
+
 	fmt.Println(dmp.DiffPrettyText(diffs))
 }

```

## Example
```go
import (
	"fmt"
	"time"

	"github.com/ogios/go-diffcontext/diffline"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const code1 = `
const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)
`

const code2 = `
const (
	text2 = "Lorem sit amet."
)
`
package main


func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(code1, code2, true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupSemantic(diffs)

	dc := diffline.New()
	dc.AddDiffs(diffs)
	fmt.Printf("dc.Lines: %v\n", dc.Lines)
	fmt.Printf("dc.GetBefore(): %v\n", dc.GetBefore())
	fmt.Printf("dc.GetAfter(): %v\n", dc.GetAfter())
	fmt.Printf("dc.GetMixed(): %v\n", dc.GetMixed())
}
```

