# go-stream
Lazy stream to work with slices wich will be calculated only when terminal function called or during the wrapping
## Creation
`New[V any](slice ...V) Stream[V]` - to create stream instance

`Wrap[V any, NV any](s Stream[V], wrapValue WrapValueFunc[V, NV]) Stream[NV]` - executing stream, converts it's values to the new type and creates stream of new type 

`Join[V any, W any, VW any, K comparable](sLeft Stream[V], extractKeyLeft ExtractComparableKeyFunc[V, K], sRight Stream[W], extractKeyRight ExtractComparableKeyFunc[W, K], merge MergeValuesFunc[V, W, VW]) Stream[VW]` - joins to streams and creates stream of new type based on mapping

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




