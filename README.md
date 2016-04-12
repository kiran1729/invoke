# invoke

invoke package is a utility wrapper to the reflect.Call function in golang.

CallFunc provides a mechanism to call a named method on any object(struct).
The input parameters are taken as interfaces but constructed at runtime by
looking up the definition of the target method. This makes it easy to invoke
any function on any object using reflection. It can be used to write text/JSON
based function calls and this library to call those functions at runtime.

CalFuncWithRaw can be used to invoke any named function on an object(struct)
in golang where the inputs are partially unmarshaled JSON raw messages.
See the tests for how a function can be invoked by supplying the params as JSON
and not having to unmarshal the parameters fully before being able to call
a method on a struct.

These functions can be used to build JSON text file driven calls to methods of
structs. We use these in a framework to call test functions on structs without
having to write tests in code.

