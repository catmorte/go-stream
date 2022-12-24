# go-stream
Lazy stream to work with slices wich will be calculated only when terminal function called or during the wrapping
## Creation
New
```
import "github.com/catmorte/go-streams/pkg/stream"
...
array := []SomeType1{...} 
...
stream := stream.New(array)
```
Mapped from another
```
import "github.com/catmorte/go-streams/pkg/stream"
...
array := []SomeType1{...} 
...
streamType1 := stream.New(array)
streamType2 := stream.Wrap(streamType1, func(i int, v SomeType1) []SomeType2 {
  return []SomeType2{...}
})
```
## Methods
### Non-terminal methods
`Sort(sortValues SortFunc[V]) Stream[V]` - sort values with function

`Filter(checkValues FilterFunc[V]) Stream[V]` - filter values with function

`Peek(peekValues PeekFunc[V]) Stream[V]` - peek values with function

`Limit(int) Stream[V]` - limit values to number

`Skip(int) Stream[V]` - skip values from number

`Distinct(compareValues EqFunc[V]) Stream[V]` - remove duplicates using funcition to compare 

`DistinctByKey(compareValues ExtractKeyFunc[V]) Stream[V]` - remove duplicates using funcition to create key 

`Expand(remap ExpandFunc[V]) Stream[V]` - expand value element to a slice of values

`Reverse() Stream[V]` - reverse slice order

### Terminal methods

`First() (V, bool)` - get first value (true if exists)

`FirstBy(checkValues FilterFunc[V]) (int, V, bool)` - get first value that is satisfying condition func returns index, value and if exists (true if exists)

`Last() (V, bool)` - get last value (true if exists)

`LastBy(checkValues FilterFunc[V]) (int, V, bool)` - get last value that is satisfying condition func returns index, value and if exists (true if exists)

`Count() int` - get amount of values

`AllMatch(checkValues FilterFunc[V]) bool` - check if all matches condition function

`AnyMatch(checkValues FilterFunc[V]) bool` - check if any matches condition function

`NoneMatch(checkValues FilterFunc[V]) bool` - check if none matches condition function

`ForEach(do DoFunc[V]) error` - do for each synchronously

`ForEachAsync(do DoFunc[V]) error` - do for each asynchronously

`ForEachChunk(chunkSize int, do DoChunkFunc[V]) error` - do for each chunk of size synchronously

`ForEachChunkAsync(chunkSize int, do DoChunkFunc[V]) error` - do for each chunk of size asynchronously
    
`Get() []V` - get values




