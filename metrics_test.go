package hystrix

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNumberIncrementAndSum(t *testing.T) {
	number := newNumber()
	var i int64
	now := time.Now().Unix()

	for i = 0; i < 20; i++ {
		number.increment(now)
	}
	sum := number.sum()
	var ret float64 = 20
	assert.Equal(t, ret, sum)

	number.clear()

	for i = 0; i < 20; i++ {
		number.increment(now - i)
	}
	sum = number.sum()
	ret = 10
	assert.Equal(t, ret, sum)
}

func TestMetricSuccess(t *testing.T) {
	metric := newMetrics()
	metric.metricSuccess(time.Now().Unix())

	sum := metric.success.sum()
	var ret float64 = 1
	assert.Equal(t, ret, sum)
}

func TestMetricFail(t *testing.T) {
	metric := newMetrics()
	metric.metricFail(time.Now().Unix())

	sum := metric.fail.sum()
	var ret float64 = 1
	assert.Equal(t, ret, sum)
}

func TestErrorPercent(t *testing.T) {
	metric := newMetrics()
	now := time.Now().Unix()

	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			metric.metricSuccess(now)
		} else {
			metric.metricFail(now)
		}
	}

	percent := metric.errorPercent()

	assert.Equal(t, 50, percent)
}
