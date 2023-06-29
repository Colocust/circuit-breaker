package main

import (
	"errors"
	"fmt"
	"github.com/Colocust/hystrix"
)

func main() {
	h := hystrix.Get()
	h.ConfigureHystrix("test", &hystrix.Circuit{
		RequestVolumeThreshold: 5,
		ErrorPercentThreshold:  20,
		SleepWindow:            1000,
	})

	for i := 0; i < 10; i++ {
		_ = h.Do("test", func() error {
			return test()
		}, nil)
	}
}

func test() error {
	fmt.Println("s")
	return errors.New("ss")
}
