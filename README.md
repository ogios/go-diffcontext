# go-diffcontext

> Based on [`sergi/go-diff`](https://github.com/sergi/go-diff)

transform from 
```go
// Diff represents one diff operation
type Diff struct {
	Type Operation
	Text string
}
```
to
```go
type DiffLine struct {
	Before []byte
	After  []byte
	State  diffmatchpatch.Operation
}
```

and able to get both content before and after change by `GetBefore` and `GetAfter`
