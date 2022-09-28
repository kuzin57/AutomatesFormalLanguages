package main

import "workspace/internal/shell"

func main() {
	shell := &shell.Shell{}
	err := shell.Init()
	if err != nil {
		panic("error while initializing")
	}
	shell.Run()

}
