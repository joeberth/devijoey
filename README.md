# devijoey
Deviget Golang-Challenge Transparent Cache
[Joeberth Souza](https://www.linkedin.com/in/joeberth-souza-5624a112a/)
joeberthaugusto@gmail.com

## Technical Decisions Taken

### First TODO in CODE ###
- *TODO: check that the price was retrieved less than "maxAge" ago!*

Before start, I created the **PriceAndTime** structure that carries the price value and the time of the last time the value was changed.

```go
type PriceAndTime struct {
	timeModified time.Time
	value        float64
}

```

To check if the price retrieved its valid in cache, a  TransparentCache funcion **IsValidCache**, this func will return true if the itemCode value retrieve is less than maxAge.

```go
func (c *TransparentCache) IsValidCache(itemCode string) bool {
	return time.Now().Sub(c.prices[itemCode].timeModified) < c.maxAge
}

```

This function is the key to solve the first TODO, but before that we need to be safe about mutual exclusion, so in **TransparentCache** struct, I added a mutex to protect the memory region area. So, before we do the reading, we use Lock to ensure that memory region doesn't have more than one thread at a time. A similar process was used to change the value on the map when the information is not valid in the cache.

### Second TODO in CODE ###

- *TODO: parallelize this, it can be optimized to not make the calls to the external service sequentially*

This second TODO needed a *errChannel* to handle with errors, and a WaitGroup to handle the goRoutines at each service call, with this we were able to return the errors and correct answers in parallel.

## TESTS

All tests passed without problems. (100%)

I added tests for a wider range of simultaneous calls, and the implementation also worked perfectly.

Curious Fact:
When making simultaneous requests for the same item, the API Might need to be called more than once. The perfect solution would be if the API were called only once.
I also added **TestRepeatedItem_ParallelizeCalls** To show that maybe there still have a space for improvement.

### Code Coverage

Using fmtcoverage lib, we got great results:
ok      fmt     0.128s  coverage: 95.2% of statements


