utils package contains some useful functions for the layout package.

for example, the function `throw` is used to throw panic when an error is not nil.
```go
func throw(err error) {
    if err != nil {
        panic(err)
    }
}
```