package plat

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

// Topic 常量定义
const (
	TopicSysTick       = "sys_tick"       // 系统 1 秒定时器
	TopicParamChanged  = "param_changed"  // 参数值变化: (index, value string)
	TopicDeviceOnline  = "device_online"  // 设备上线: (deviceIndex string)
	TopicDeviceOffline = "device_offline" // 设备离线: (deviceIndex string)
	TopicAlarm         = "alarm"          // 告警: (index, alarmType, value string)
	TopicCommand       = "command"        // 下发命令: (index, value string)
)

// EB 事件总线单例
type EB struct {
	bus EventBus.Bus
}

var (
	ebInstance *EB
	ebOnce     sync.Once
)

// GetEventBus 获取单例
func GetEventBus() *EB {
	ebOnce.Do(func() {
		ebInstance = &EB{
			bus: EventBus.New(),
		}
	})
	return ebInstance
}

// Subscribe 同步订阅
func (e *EB) Subscribe(topic string, fn interface{}) {
	e.bus.Subscribe(topic, fn)
}

// SubscribeAsync 异步订阅
func (e *EB) SubscribeAsync(topic string, fn interface{}, transactional bool) {
	e.bus.SubscribeAsync(topic, fn, transactional)
}

// Publish 发布事件
func (e *EB) Publish(topic string, args ...interface{}) {
	e.bus.Publish(topic, args...)
}
