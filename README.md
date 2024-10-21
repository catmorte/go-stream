# go-stream
A lazy stream implementation for working with slices, where computations are deferred until a terminal function is called or during wrapping.

## Creation
- **`New[V any](slice ...V) Stream[V]`**:  
  Creates a new stream instance from a slice.

- **`Wrap[V any, NV any](s Stream[V], wrapValue WrapValueFunc[V, NV]) Stream[NV]`**:  
  Executes the stream, converts its values to a new type, and returns a stream of the new type.

- **`Join[V any, W any, VW any, K comparable](sLeft Stream[V], extractKeyLeft ExtractComparableKeyFunc[V, K], sRight Stream[W], extractKeyRight ExtractComparableKeyFunc[W, K], merge MergeValuesFunc[V, W, VW]) Stream[VW]`**:  
  Joins two streams, mapping and merging values into a new stream of the specified type.

## Non-terminal Methods
- **`Sort(sortValues SortFunc[V]) Stream[V]`**:  
  Sorts stream values using the provided function.

- **`Filter(checkValues FilterFunc[V]) Stream[V]`**:  
  Filters stream values based on the provided function.

- **`Peek(peekValues PeekFunc[V]) Stream[V]`**:  
  Peeks at stream values using the provided function.

- **`Limit(limit int) Stream[V]`**:  
  Limits the stream to the specified number of values.

- **`Skip(skip int) Stream[V]`**:  
  Skips the specified number of values in the stream.

- **`Distinct(compareValues EqFunc[V]) Stream[V]`**:  
  Removes duplicates using the provided comparison function.

- **`DistinctByKey(extractKey ExtractKeyFunc[V]) Stream[V]`**:  
  Removes duplicates based on a key extraction function.

- **`Expand(remap ExpandFunc[V]) Stream[V]`**:  
  Expands each element into a slice of values.

- **`Reverse() Stream[V]`**:  
  Reverses the order of values in the stream.

## Terminal Methods
- **`Get() []V`**:  
  Returns the computed values as a slice.

- **`ForEach(do DoFunc[V]) error`**:  
  Executes an action for each value in the stream, synchronously.

- **`ForEachAsync(do DoFunc[V]) error`**:  
  Executes an action for each value in the stream, asynchronously.

- **`ForEachChunk(chunkSize int, do DoChunkFunc[V]) error`**:  
  Executes an action for each chunk of values in the stream, synchronously.

- **`ForEachChunkAsync(chunkSize int, do DoChunkFunc[V]) error`**:  
  Executes an action for each chunk of values in the stream, asynchronously.

- **`Count() int`**:  
  Returns the number of values in the stream.

- **`First() (V, bool)`**:  
  Returns the first value, or `false` if the stream is empty.

- **`FirstBy(checkValues FilterFunc[V]) (int, V, bool)`**:  
  Returns the index and first value that satisfies the condition, or `false` if no match is found.

- **`Last() (V, bool)`**:  
  Returns the last value, or `false` if the stream is empty.

- **`LastBy(checkValues FilterFunc[V]) (int, V, bool)`**:  
  Returns the index and last value that satisfies the condition, or `false` if no match is found.

- **`AllMatch(checkValues FilterFunc[V]) bool`**:  
  Returns `true` if all values match the condition.

- **`AnyMatch(checkValues FilterFunc[V]) bool`**:  
  Returns `true` if any value matches the condition.

- **`NoneMatch(checkValues FilterFunc[V]) bool`**:  
  Returns `true` if no values match the condition.


