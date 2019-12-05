package main

import (
	"fmt"

	"github.com/zmitry/go2typings"
	"github.com/zmitry/go2typings/example/types"
)

func main() {
	s := go2typings.New()
	s.Add(types.T{})
	s.Add(types.User{})

	str, _ := s.RenderToSwagger()
	fmt.Print(str)
	// if err != nil {
	// 	panic(err)
	// }
}
