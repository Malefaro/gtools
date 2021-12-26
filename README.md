# Gtools 
_________________
Generic tools for go 1.18+ 

## FT (func tools)
______
Provide func tools over iterators

Iterators for functions like `Filter`, `Map`, `Reduce`, `etc` solve 3 main problems: 
1. Does not allocate new slices (because you just iterates over provided one)
2. Iterates over slice just once (without iterators in case chaining filter -> map -> reduce  iterates 3 times [1 time in each function])
3. Iterator can be writen for any type you needed. So you can work with custom types (such a `List`, `Queue`, `etc` ) in same way

### Examples:
Filter -> Map -> Reduce:

```
iter := ft.SliceIter([]int{1,2,3,4})
result := ft.Reduce(ft.Map(ft.Filter(iter, func(t int) bool {
	return t%2 != 0
}), func(t int) int {
	return t * 2
}), func(o string, t int) string {
	return o + strconv.Itoa(t)
})
```

Zip: 

```
iter1 := ft.SliceIter([]int{1,2,3})
iter2 := ft.SliceIter([]string{"one", "two", "three")
iter3 := ft.SliceIter([]float64{1.1, 2.2, 3.3)
for _, p := range ft.Collect(ft.Zip(ft.Zip(iter1, iter2), iter3)) {
	fmt.Println(p.First.First, p.First.Second, p.Second)
}
```

Count:

```
iter := ft.SliceIter([]int{1,2,3)
result := ft.Count(iter) // 3
// with predicate
iter2 := ft.SliceIter([]int{1,2,3})
result := ft.Count(iter, func(t int) bool {
    return t%2 == 0
}) // 1

```

### This package contains:

##### Functions:
* `Filter` - wrap provided iter and return new iterator, that yields elements satisfying the given function
* `Map` - wrap provided iter and return new iterator, that yields elements obtained by applying the given function to each element of the original iterator
* `Skip` - skip N iterations of iterator
* `Chunk` - split provided iter into several slices returns iterator of slices with len less or equal `size`
* `Scan` -  same as reduce but instead of returning one result it return iterator of results at every step
* `Product` - make cartesian product of input iters.
* `Cycle` - return endless iterator that yields elements from original iter
* `Zip` - create a new iterator over provided 2. This iterator yields pairs of each iterator elements. It ends when one of iter ends (you can combine it if you need zip more than 2 iters: `Zip(Zip(iter1, iter2), iter3)`)
* `Enumerate` - returns an iterator of the original slice elements with numbering

##### Consumers:
* `Collect` - consumes iterator and return slice of its elements
* `CollectInto` -  consumes the iterator and fill provided arguments
* `CollectR` - same as `CollectInto` but this func returns value of provided (during call) type. This func uses reflect.
* `Any`- consumes iter and returns true if any element of iter returns true on predicate func call on it
* `All` - consumes iter and returns true if all elements returns true on predicate func call on it
* `Reduce` - consumes iter, calls provided func on every element of iterator, accumulating result (has optional argument for inital value)
* `Count` - consumes iter and count its elements (has optional argument to call only specific values)
* `Sum` - consumes elements and return sum of it (can be used only on iterators with numeric types (such as int. float, complex)
* `ForEach` - consumes iterator and apply provided function on it 
* `IntoChannel` - return channel that yields iterator elements (has optional arg context.Context)
* `Max` - consumes iter and return max element find in it or nil if no such element
* `Min` - same as `Max` but return min element
* `Contains` - return true if iterator contains provided element

###### Iterator consctructors:
* `SliceIter` - iterator over slice
* `MapIter` - iterator over map (this iterator spawn goroutine to read from map) use this if you have huge size map
* `MapIterOverSlice` - iterator over map (this iterator creating `SliceIter` with all key-value pairs)

