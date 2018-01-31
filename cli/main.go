package main

import (
	"fmt"

	"github.com/open-pomodoro/go-openpomodoro"
)

func main() {
	opc, err := openpomodoro.NewClient("/tmp/.pomodoro")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", opc)

	history, err := opc.History()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", history.Count())
}
