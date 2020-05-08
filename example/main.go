package main

import (
	"encoding/json"
	"fmt"

	"github.com/zmitry/go2typings"
	"github.com/zmitry/go2typings/example/types"
)

func main() {
	s := go2typings.New()
	s.Add(types.T{})
	s.Add(types.User{})

	str := s.RenderToSwagger()
	res, _ := json.MarshalIndent(str, "", " ")
	fmt.Print(string(res))
	// if err != nil {
	// 	panic(err)
	// }
}
