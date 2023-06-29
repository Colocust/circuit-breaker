package hystrix

import (
	"fmt"
	"testing"
)

func BenchmarkHystrix(b *testing.B) {
	h := Get()
	h.ConfigureHystrix("test", &Circuit{
		RequestVolumeThreshold: 5,
		ErrorPercentThreshold:  20,
		SleepWindow:            1000,
	})

	for n := 0; n < b.N; n++ {
		_ = h.Do("test", func() error {
			return nil
		}, func(err error) error {
			fmt.Println("err:", err)
			return nil
		})
	}
}
