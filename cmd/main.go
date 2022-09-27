package main

import (
	"workspace/internal/shell"
)

// func F[T Stringer](t T) string {
// 	return t.String()
// }

// type Custom struct {
// 	S string
// }

// func (c Custom) String() string {
// 	return c.S
// }

func main() {
	// c := Custom{S: "haha"}
	// fmt.Println(F(c))
	shell := &shell.Shell{}
	err := shell.Init()
	if err != nil {
		panic("error while initializing")
	}
	shell.Run()
}
