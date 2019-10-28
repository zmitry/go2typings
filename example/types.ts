/* tslint:disable */
/* eslint-disable */
export namespace types {
  //github.com/zmitry/go2typings/example/types.M
  export interface M {
    Username2: string;
  }
  //.
  export interface NameT {}
  //github.com/zmitry/go2typings/example/types.T
  export interface T extends M {
    name: Array<Record<string, types.NameT>> | null;
    lastname: Array<Record<string, string>> | null;
    surname: Array<Record<string, Array<types.M | null> | null>> | null;
  }
  //github.com/zmitry/go2typings/example/types.UserTag
  export interface UserTag {
    tag: string;
  }
  //github.com/zmitry/go2typings/example/types.User
  export interface User {
    firstname: string;
    secondName: string;
    tags: Array<types.UserTag> | null;
  }
}

