package main

const code1 = `
package main

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
}
`

const code2 = `
package main

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem sit amet."
)

func main() {
	diffs := dmp.DiffMain(text1, text2, false)

	dmp := diffmatchpatch.New()

	fmt.Println(dmp.DiffPrettyText(diffs))
}
`

// const code1 = `
// 	dmp := diffmatchpatch.New()
//
// 	diffs := dmp.DiffMain(text1, text2, false)
//
// 	fmt.Println(dmp.DiffPrettyText(diffs))
// `
//
// const code2 = `
// 	diffs := dmp.DiffMain(text1, text2, false)
//
// 	dmp := diffmatchpatch.New()
//
// 	fmt.Println(dmp.DiffPrettyText(diffs))
// `
