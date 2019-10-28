# Golang structs to typescript typings convertor

Here is the cases we handle with this lib

```golang
package userapi
type M struct {
	Username string `json:"Username2"`
 }
 type T struct {
 	M
 	Name []map[string]struct {
 		test string
 	} `json:"name"`

  // Array<Record<string, string>>
  lastname []map[string]string `json:"lastname"`
  surname  []map[string][]*M   `json:"surname"`
}
```

output

```ts
export namespace userapi {
  //userapi.Root
  export interface Root {
    info: string;
  }
  //userapi.M
  export interface M extends Root {
    Username2: string;
  }
  //.
  export interface NameT {}
  //userapi.T
  export interface T extends M {
    name: Array<Record<string, main.NameT>> | null;
    Name2: main.NameT;
    lastname: Array<Record<string, string>> | null;
    surname: Array<Record<string, Array<main.M | null> | null>> | null;
  }
}
```

to see working example go to /example

# How to setup

create go file with the following code

```golang
package main

import (
  "github.com/zmitry/go2ts"
  // you can use your own
	"github.com/zmitry/go2ts/example/types"
)

type Root struct {
	User types.User
	T    types.T
}

func main() {
	s := go2ts.New(&go2ts.Options{})
	s.Add(types.T{})
	s.Add(types.User{})

	err := s.GenerateFile("./test.ts")
	if err != nil {
		panic(err)
	}
}
```

# Custom tags

we support custom tag `ts` it supports the following syntax

```
type M struct {
	Username string `json:"Username2" ts:"string,optional"`
}
```

tsTag type

```
tsTag[0] = "string"|"date"|"-"
tsTag[1] = "optional"|"no-null"|"null"
```

see field.go for more info

# TODO:

- add tests
- add customization for intendation and output format
