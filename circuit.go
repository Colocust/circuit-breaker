package hystrix

import (
	"sync"
	"time"
)

type Circuit struct {
	RequestVolumeThreshold int // 达到这个请求数量后才去判断是否要开启熔断
	ErrorPercentThreshold  int // 请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	SleepWindow            int // 熔断器被打开后 SleepWindow 的时间就是控制过多久后去尝试服务是否可用了 单位为毫秒

	open         bool
	lastOpenTime int64 // 单位ms
	mutex        sync.RWMutex

	metric *metrics
}

type eventType int

const (
	defaultRequestVolumeThreshold = 20
	defaultErrorPercentThreshold  = 50
	defaultSleepWindow            = 5000
)

const (
	successEvent     eventType = iota // 成功
	circuitOpenEvent                  // 熔断器被打开
	failureEvent                      // 执行业务逻辑错误
	fallbackSuccessEvent
	fallbackFailEvent
)

// 获取一个默认熔断器
func getDefaultCircuit() *Circuit {
	circuit := new(Circuit)
	circuit.completeConfigure()
	return circuit
}

func (circuit *Circuit) allowRequest() bool {
	return !circuit.isOpen()
}

func (circuit *Circuit) reportEvent(events []eventType) {
	now := time.Now().Unix()

	for _, event := range events {
		switch event {
		case successEvent:
			circuit.metric.metricSuccess(now)
		case failureEvent:
			circuit.metric.metricFail(now)
		}
	}

	if !circuit.isHealthy() {
		circuit.setOpen()
	}
}

// 熔断器是否打开
func (circuit *Circuit) isOpen() bool {
	circuit.mutex.RLock()
	o := circuit.open
	circuit.mutex.RUnlock()

	if !o {
		return false
	}

	if circuit.lastOpenTime+int64(circuit.SleepWindow) < time.Now().UnixMilli() {
		circuit.setClose()
		return false
	}
	return true
}

func (circuit *Circuit) setClose() {
	circuit.mutex.Lock()
	defer circuit.mutex.Unlock()

	if !circuit.open {
		return
	}

	circuit.open = false
	circuit.metric.clear()
}

func (circuit *Circuit) setOpen() {
	circuit.mutex.Lock()
	defer circuit.mutex.Unlock()

	if circuit.open {
		return
	}

	circuit.open = true
	circuit.lastOpenTime = time.Now().UnixMilli()
}

func (circuit *Circuit) isHealthy() bool {
	if int(circuit.metric.totalRequest()) < circuit.RequestVolumeThreshold {
		return true
	}
	return circuit.metric.errorPercent() < circuit.ErrorPercentThreshold
}

func (circuit *Circuit) completeConfigure() {
	if circuit.RequestVolumeThreshold == 0 {
		circuit.RequestVolumeThreshold = defaultRequestVolumeThreshold
	}
	if circuit.ErrorPercentThreshold == 0 || circuit.ErrorPercentThreshold > 100 {
		circuit.ErrorPercentThreshold = defaultErrorPercentThreshold
	}
	if circuit.SleepWindow == 0 {
		circuit.SleepWindow = defaultSleepWindow
	}

	if circuit.metric == nil {
		circuit.metric = newMetrics()
	}
}
