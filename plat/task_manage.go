package plat

import (
	"context"
	"fmt"
	"sync"

	"SH_SCADA/common/model"

	"gorm.io/gorm"
)

// TaskManager 采集任务管理器（单例）
type TaskManager struct {
	db *gorm.DB

	// 任务静态数据
	tasks   []*model.Task
	taskMap map[string]*model.Task

	// 已注册的驱动 name → Driver
	drivers map[string]Driver

	// 运行中的任务 ctx 取消函数
	cancels map[string]context.CancelFunc

	mu sync.Mutex
}

var (
	taskInstance *TaskManager
	taskOnce     sync.Once
)

// GetTaskManager 获取单例
func GetTaskManager() *TaskManager {
	taskOnce.Do(func() {
		taskInstance = &TaskManager{
			drivers: make(map[string]Driver),
			cancels: make(map[string]context.CancelFunc),
		}
	})
	return taskInstance
}

// RegisterDriver 注册采集驱动（main 中调用）
func (tm *TaskManager) RegisterDriver(d Driver) {
	tm.drivers[d.Info()] = d
}

// Init 从数据库加载任务（仅启动时调用一次）
func (tm *TaskManager) Init(db *gorm.DB) error {
	tm.db = db

	var raw []model.Task
	if err := db.Find(&raw).Error; err != nil {
		return err
	}

	tm.tasks = make([]*model.Task, 0, len(raw))
	tm.taskMap = make(map[string]*model.Task, len(raw))
	for i := range raw {
		t := &raw[i]
		tm.tasks = append(tm.tasks, t)
		tm.taskMap[t.Index] = t
	}

	return nil
}

// StartAll 启动所有启用的任务（main 中调用，gorm.DB 初始化之后）
func (tm *TaskManager) StartAll(mgr *DevicesManager) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, t := range tm.tasks {
		if !t.Enable {
			continue
		}
		tm.startLocked(t, mgr)
	}
}

// StopAll 停止所有任务（用于优雅退出）
func (tm *TaskManager) StopAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for index, cancel := range tm.cancels {
		cancel()
		delete(tm.cancels, index)
	}
}

// GetTask 获取单个任务
func (tm *TaskManager) GetTask(index string) *model.Task {
	return tm.taskMap[index]
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() []*model.Task {
	return tm.tasks
}

// startLocked 内部方法，调用前需持有 mu
func (tm *TaskManager) startLocked(t *model.Task, mgr *DevicesManager) {
	d, ok := tm.drivers[t.Driver]
	if !ok {
		fmt.Printf("[TaskManager] 未找到驱动 %s (任务: %s)\n", t.Driver, t.Index)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	tm.cancels[t.Index] = cancel

	go func(task *model.Task) {
		defer cancel()
		fmt.Printf("[TaskManager] 启动采集 %s (驱动: %s)\n", task.Index, task.Driver)
		if err := d.Start(task, mgr, ctx); err != nil {
			fmt.Printf("[TaskManager] 采集异常 %s: %v\n", task.Index, err)
		}
	}(t)
}
