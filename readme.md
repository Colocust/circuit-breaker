## 快速开始

获取hystrix对象：

```golang
h := hystrix.Get()
h.ConfigureHystrix("test", &hystrix.Circuit{
    RequestVolumeThreshold: 5,
    ErrorPercentThreshold:  20,
    SleepWindow:            1000,
})
```

执行代码:

```golang
h.Do(context.Background(), "test", func () error {
    // do something
    return nil
}, func (err error) error {
    // handle err
    fmt.Println(err)
})
```

## 配置

| 配置名                    | 解释                                 | 默认值  |
|------------------------|------------------------------------|------|
| RequestVolumeThreshold | 在10s的采样时间内，请求次数打到这个阈值后，才开始判断是否开启熔断 | 20   |
| ErrorPercentThreshold  | 错误率阈值。超过后，会启用熔断和降级处理。              | 50   |
| SleepWindow            | 休眠窗口期。熔断后等待试探（是否关闭熔断器）的时间。         | 5000 |