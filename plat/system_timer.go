package plat

import (
	"sync"
	"time"
)

// SysTimer 系统 1 秒定时器，通过事件总线发布 TopicSysTick
type SysTimer struct {
	stop chan struct{}
	eb   *EB
}

var (
	sysTimerInstance *SysTimer
	sysTimerOnce     sync.Once
)

// GetSysTimer 获取单例
func GetSysTimer() *SysTimer {
	sysTimerOnce.Do(func() {
		sysTimerInstance = &SysTimer{
			stop: make(chan struct{}),
			eb:   GetEventBus(),
		}
	})
	return sysTimerInstance
}

// Start 启动定时器，每秒通过事件总线发布 TopicSysTick
func (st *SysTimer) Start() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-st.stop:
				ticker.Stop()
				return
			case <-ticker.C:
				st.eb.Publish(TopicSysTick)
			}
		}
	}()
}

// Stop 停止定时器
func (st *SysTimer) Stop() {
	close(st.stop)
}
