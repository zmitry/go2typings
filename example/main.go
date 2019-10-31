package main

import (
	"github.com/zmitry/go2typings"
	"github.com/zmitry/go2typings/example/types"
)

func main() {
	s := go2typings.New()
	s.Add(types.T{})
	// s.Add(types.User{})

	err := s.GenerateFile("./types.ts")
	if err != nil {
		panic(err)
	}
}
