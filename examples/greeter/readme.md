Greeter Example
================
This example demonstrates a rather simple input based "Hello World" example. The name provided in the input is displayed with greeting.
But more than anything, it was created to demonstrated the creation of the Components, Views and Apps, It also demonstrates the new 
event system for Gu which abstracts away all details for a more Go centric approach.


## Building
Simple navigate your terminal to this directory and run the following command.

```go
gopherjs build -o assets/greeter.js
```

### Running
To have the exampe running, simple start the gopherjs server `gopherjs serve` and navigate to `localhost:8080/github.com/gu-io/examples/greeter`.
