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

package invoke

import (
  "encoding/json"
  "fmt"
  "strings"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"

  "github.com/zerostackinc/customtypes"
)

type myStruct struct {
  member int
}

type SampleStruct struct {
  Params []customtypes.RawMessage `json:"params"`
}

func (m *myStruct) ExampleFunc(intArg int, stringArg string,
  sliceArg []int, dur time.Duration) (int, int64, error) {

  return 100, dur.Nanoseconds(), nil
}

func TestCallFunc(t *testing.T) {
  my := myStruct{}
  var myInt int
  myString := "examplestring"
  mySlice := []int{1, 2, 3}
  myDur := 13 * time.Second
  results := CallFunc(&my, "ExampleFunc", myInt, myString, mySlice, myDur)
  assert.NotNil(t, results)

  assert.Equal(t, len(results), 3, fmt.Sprintf("results = %#v", results))

  resInt, ok := results[0].Interface().(int)
  assert.True(t, ok)
  assert.Equal(t, resInt, 100)

  resInt64, ok := results[1].Interface().(int64)
  assert.True(t, ok)
  assert.Equal(t, resInt64, int64(13000000000))

  assert.Equal(t, results[2].Interface(), nil)

  // Check that it does recover from panic when using wrong params.
  // TODO: when we fix the type checking TODO in CallFunc then it will not
  // panic so this needs to incite panic in a different way.
  results = CallFunc(&my, "ExampleFunc", myString, myInt, mySlice, myDur)

  assert.Equal(t, len(results), 1)

  resErr, ok := results[0].Interface().(error)
  assert.True(t, ok)
  assert.True(t, strings.Contains(resErr.Error(), "recovered from panic"))
}

func TestCallFuncRaw(t *testing.T) {
  sampleJSON := []byte(`
  {
    "params" : [ 2, "examplestring", [5, 6, 7, 8], 10 ]
  }
`)

  var sampleInput SampleStruct

  err := json.Unmarshal(sampleJSON, &sampleInput)

  assert.NoError(t, err)
  assert.Equal(t, len(sampleInput.Params), 4)

  my := myStruct{}

  results := CallFuncWithRaw(&my, "ExampleFunc", sampleInput.Params)
  assert.NoError(t, err)
  assert.NotNil(t, results)
  assert.Equal(t, len(results), 3)

  resInt, ok := results[0].Interface().(int)
  assert.True(t, ok)
  assert.Equal(t, resInt, 100)

  resInt64, ok := results[1].Interface().(int64)
  assert.True(t, ok)
  assert.Equal(t, resInt64, int64(10))

  assert.Equal(t, results[2].Interface(), nil)
}
