package hystrix

import (
	"math"
	"sync"
	"time"
)

type (
	metrics struct {
		total   *number
		success *number
		fail    *number
	}
	number struct {
		buckets [10]*bucket
		mutex   sync.RWMutex
	}
	bucket struct {
		timestamp int64
		value     float64
	}
)

func newMetrics() *metrics {
	return &metrics{
		total:   newNumber(),
		success: newNumber(),
		fail:    newNumber(),
	}
}

func newNumber() *number {
	return &number{
		buckets: newBuckets(),
		mutex:   sync.RWMutex{},
	}
}

func newBuckets() (buckets [10]*bucket) {
	for i := 0; i < 10; i++ {
		buckets[i] = new(bucket)
	}
	return
}

func (metric *metrics) metricSuccess(now int64) {
	metric.total.increment(now)
	metric.success.increment(now)
}

func (metric *metrics) metricFail(now int64) {
	metric.total.increment(now)
	metric.fail.increment(now)
}

func (metric *metrics) totalRequest() float64 {
	return metric.total.sum()
}

func (metric *metrics) clear() {
	metric.total.clear()
	metric.success.clear()
	metric.fail.clear()
}

func (metric *metrics) errorPercent() int {
	total, fail := metric.totalRequest(), metric.fail.sum()
	percent := (fail / total) * 100
	return int(math.Floor(percent + 0.5)) // 四舍五入
}

func (number *number) increment(now int64) {
	number.mutex.Lock()
	defer number.mutex.Unlock()

	index := now % 10

	if now < number.buckets[index].timestamp {
		return
	}

	if number.buckets[index].timestamp != now {
		number.buckets[index].value = 0
	}
	number.buckets[index].timestamp = now
	number.buckets[index].value++
}

func (number *number) clear() {
	number.mutex.Lock()
	defer number.mutex.Unlock()

	number.buckets = newBuckets()
}

func (number *number) sum() (sum float64) {
	number.mutex.RLock()
	defer number.mutex.RUnlock()

	now := time.Now().Unix()
	for _, ele := range number.buckets {
		if ele.timestamp <= now-10 {
			continue
		}
		sum += ele.value
	}
	return
}
