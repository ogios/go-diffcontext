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
