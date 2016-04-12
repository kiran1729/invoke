// Copyright 2014 ZeroStack, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// invoke package provides functions to invoke methods on objects(structs) by
// passing a pointer to an object, the name of the method and the params.
// The params can also be passed as customtypes.RawMessage slice which will
// be unmarshaled into specific types as expected by the method.

package invoke

import (
  "encoding/json"
  "fmt"
  "reflect"
  "runtime/debug"

  "github.com/zerostackinc/customtypes"
)

// CallFunc invokes the method with name on object "any" using the params.
func CallFunc(any interface{}, name string, params ...interface{}) (
  results []reflect.Value) {

  defer func() {
    if r := recover(); r != nil {
      errP := fmt.Errorf("invoke: recovered from panic in CallFunc :: %v %s", r,
        debug.Stack())
      results = []reflect.Value{reflect.ValueOf(errP)}
    }
  }()

  if any == nil {
    errN := fmt.Errorf("invoke: nil interface passed for Validate")
    return []reflect.Value{reflect.ValueOf(errN)}
  }

  typ := reflect.TypeOf(any)
  val := reflect.ValueOf(any)
  if typ.Kind() != reflect.Ptr {
    errT := fmt.Errorf("invoke: input interface is not a ptr : %v", typ.Kind())
    return []reflect.Value{reflect.ValueOf(errT)}
  }

  if typ.Elem().Kind() != reflect.Struct {
    errK := fmt.Errorf("invoke: input interface %v is not a ptr to struct %v",
      typ.Kind(), typ.Elem().Kind())
    return []reflect.Value{reflect.ValueOf(errK)}
  }

  _, found := typ.MethodByName(name)
  if !found {
    errM := fmt.Errorf("invoke: could not find method %s for type %v "+
      "num_methods=%d", name, typ, val.NumMethod())
    return []reflect.Value{reflect.ValueOf(errM)}
  }

  method := val.MethodByName(name)

  if len(params) != method.Type().NumIn() {
    errP := fmt.Errorf("invoke: mismatch in number of params %d and func "+
      "inputs %d", len(params), method.Type().NumIn())
    return []reflect.Value{reflect.ValueOf(errP)}
  }

  // TODO: do type checking of params to CallFunc does not panic due to
  // mismatched types.

  in := make([]reflect.Value, len(params))
  for k, param := range params {
    in[k] = reflect.ValueOf(param)
  }
  results = reflect.ValueOf(any).MethodByName(name).Call(in)
  return results
}

// CallFuncWithRaw invokes the method with name on object "any" using the params
// RawMessage slice that is unmarshaled using the methods own definition of
// what each param type is.
func CallFuncWithRaw(any interface{}, name string,
  params []customtypes.RawMessage) (results []reflect.Value) {

  defer func() {
    if r := recover(); r != nil {
      errP := fmt.Errorf("invoke: recovered from panic in CallFuncWithRaw :: "+
        "%v %s", r, debug.Stack())
      results = []reflect.Value{reflect.ValueOf(errP)}
    }
  }()

  if any == nil {
    errN := fmt.Errorf("invoke: nil interface passed for Validate")
    return []reflect.Value{reflect.ValueOf(errN)}
  }

  typ := reflect.TypeOf(any)
  val := reflect.ValueOf(any)
  if typ.Kind() != reflect.Ptr {
    errT := fmt.Errorf("invoke: input interface is not a ptr : %v", typ.Kind())
    return []reflect.Value{reflect.ValueOf(errT)}
  }

  if typ.Elem().Kind() != reflect.Struct {
    errK := fmt.Errorf("invoke: input interface %v is not a ptr to struct %v",
      typ.Kind(), typ.Elem().Kind())
    return []reflect.Value{reflect.ValueOf(errK)}
  }

  _, found := typ.MethodByName(name)
  if !found {
    errM := fmt.Errorf("invoke: could not find method %s for type %v "+
      "num_methods=%d", name, typ, val.NumMethod())
    return []reflect.Value{reflect.ValueOf(errM)}
  }

  method := val.MethodByName(name)

  if len(params) != method.Type().NumIn() {
    errP := fmt.Errorf("invoke: mismatch in number of params %d and func "+
      "inputs %d", len(params), method.Type().NumIn())
    return []reflect.Value{reflect.ValueOf(errP)}
  }

  // TODO: Have to check if function signature is with pointer type
  // versus value type and what the code below has to do differently.

  in := make([]reflect.Value, len(params))
  for k, param := range params {
    inp := reflect.ValueOf(any).MethodByName(name).Type().In(k)
    inVar := reflect.New(inp)
    // unmarshal the RawMessage into the newly created var of the appropriate
    // type allocated using reflect.New()
    err := json.Unmarshal(param, inVar.Interface())
    if err != nil {
      errU := fmt.Errorf("invoke: error unmarshaling param[%d] %v for func %s"+
        " :: %v", k, param.String(), name, err)
      return []reflect.Value{reflect.ValueOf(errU)}
    }
    in[k] = inVar.Elem()
  }
  result := reflect.ValueOf(any).MethodByName(name).Call(in)
  return result
}
