package hystrix

import (
	"errors"
	"sync"
)

type (
	Hystrix struct {
		pool     sync.Pool // 获取每一个command
		circuits sync.Map  // 存储所有熔断器策略
	}

	runFunc      func() error
	fallbackFunc func(err error) error
)

var (
	ErrorRunFuncNil = errors.New("run func nil")
)

var (
	_hystrix *Hystrix
)

func init() {
	_hystrix = new(Hystrix)
	_hystrix.pool.New = func() interface{} {
		return _hystrix.allocateCommand()
	}
}

func Get() *Hystrix {
	return _hystrix
}

// ConfigureHystrix 配置熔断器
func (hystrix *Hystrix) ConfigureHystrix(name string, circuit *Circuit) {
	circuit.completeConfigure()
	hystrix.circuits.Store(name, circuit)
}

func (hystrix *Hystrix) Do(name string, run runFunc, fallback fallbackFunc) error {
	if run == nil {
		return ErrorRunFuncNil
	}

	cmd := hystrix.pool.Get().(*command)
	defer hystrix.pool.Put(cmd)

	cmd.circuit = hystrix.getCircuit(name)
	cmd.run, cmd.fallback = run, fallback
	cmd.events = make([]eventType, 0)

	cmd.do()
	return nil
}

func (hystrix *Hystrix) allocateCommand() *command {
	return new(command)
}

// 获得一个熔断器
func (hystrix *Hystrix) getCircuit(name string) *Circuit {
	if value, ok := hystrix.circuits.Load(name); ok {
		return value.(*Circuit)
	}
	hystrix.ConfigureHystrix(name, getDefaultCircuit())
	value, _ := hystrix.circuits.Load(name)
	return value.(*Circuit)
}
