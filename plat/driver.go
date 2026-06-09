package plat

import (
	"context"

	"SH_SCADA/common/model"
)

// Driver 采集驱动接口（定义在平台层，驱动包实现此接口）
type Driver interface {
	// Info 返回驱动标识，与 Task.Driver 字段匹配，内容由驱动自行决定
	Info() string
	// Start 启动采集，ctx 取消即停止
	Start(task *model.Task, mgr *DevicesManager, ctx context.Context) error
}
