package sample1

import (
	"fmt"
	"sync"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]PriceAndTime
	mutex              sync.Mutex
}

// PriceAndTime its a schema to represent when the price was modified and the Price value.
type PriceAndTime struct {
	timeModified time.Time
	value        float64
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]PriceAndTime{},
	}
}

// IsValidCache check if price ages is younger than MaxAge and return true if its can get cached data.
func (c *TransparentCache) IsValidCache(itemCode string) bool {
	return time.Now().Sub(c.prices[itemCode].timeModified) < c.maxAge
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	price, ok := c.prices[itemCode]
	if ok && c.IsValidCache(itemCode) {
		return price.value, nil
	} else if ok && !c.IsValidCache(itemCode) {
		// If its expired, we do not need the itemCode in Map anymore.
		delete(c.prices, itemCode)
	}
	priceNew, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v, Item code: %v", err.Error(), itemCode)
	}
	c.prices[itemCode] = PriceAndTime{value: priceNew, timeModified: time.Now()}
	return priceNew, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := make([]float64, len(itemCodes))
	errChannel := make(chan error, len(itemCodes))
	wg := sync.WaitGroup{}
	for i, itemCode := range itemCodes {
		// TODO: parallelize this, it can be optimized to not make the calls to the external service sequentially
		wg.Add(1)
		go func(i int, itemCode string) {
			defer wg.Done()
			price, err := c.GetPriceFor(itemCode)
			if err != nil {
				errChannel <- err
			}
			results[i] = price
		}(i, itemCode)
	}
	wg.Wait()
	close(errChannel)
	for e := range errChannel {
		return nil, e
	}

	return results, nil
}
