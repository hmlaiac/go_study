# A. Initial go project

1. Create directory and redirect to it

```shell
mkdir generics
```

```shell
cd "C:\github\go_study\Chapter6 Generics\generics"
```

2. Create go.mod file

```shell
go mod init example/generics
```

# B. Code Analysis

## B1. Non-Generic Sum

map is a very common data structure in programming languages. It helps us to search and find specific value according to
key quickly. Since the function call `map`, a hash table is created in memory. We can use loop to calculate the sum of
numbers.

`for _, v := range m`: `_` is a blank identifier, which is used to ignore the index of the map. You can replace `_`
with `i` and print out it. It will return `[0,1]` in this example.

```go
// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
var s int64
for _, v := range m {
s += v
}
return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
var s float64
for _, v := range m {
s += v
}
return s
}
```

Main function:

```go
func main() {
// Initialize a map for the integer values
ints := map[string]int64{
"first":  34,
"second": 12,
}

// Initialize a map for the float values
floats := map[string]float64{
"first":  35.98,
"second": 26.99,
}

fmt.Printf("Non-Generic Sums: %v and %v\n",
SumInts(ints),
SumFloats(floats))
}
```

Result:

```text
Non-Generic Sums: 46 and 62.97
```

## B2. Generic Sum

### B2.1 Generic function

```go
// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
var s V
for _, v := range m {
s += v
}
return s
}
```

### B2.2 Parameters

- `V`: `int64 | float64` means `V` can be either `int64` or `float64`
- `K`: `comparable` means `K` can be compared with `==` and `!=`

Q&A:
Why `K` is not `string`, but `comparable`?

-> It is because go requires the map key to be comparable. This is a <mark>constraints</mark> to key to make sure all
keys are unique.

### B2.3 Print out the results

With type arguments:

```go
fmt.Printf("Generic Sums: %v and %v\n",
SumIntsOrFloats[string, int64](ints),
SumIntsOrFloats[string, float64](floats))
```

No type arguments, type parameters inferred:

```go
fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
SumIntsOrFloats(ints),
SumIntsOrFloats(floats))
```

### B2.4 Result

```text
Generic Sums: 46 and 62.97
Generic Sums, type parameters inferred: 46 and 62.97
```

## B3. Type constraints

Declare generic type:

```go
type Number interface {
int64 | float64
}
```

Build generic function:

```go
// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SumNumbers[K comparable, V Number](m map[K]V) V {
var s V
for _, v := range m {
s += v
}
return s
}
```

# C. Conclusion

Although all functions output same results, the generic function is more flexible and can be used in more scenarios. It
is very hard to explain the generic programming in few words. Let me give you a photo to illustrate the idea of generic
programming.

![img.png](img.png)