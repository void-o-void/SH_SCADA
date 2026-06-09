package plat

import (
	"strings"
	"sync"

	"SH_SCADA/common/model"

	"gorm.io/gorm"
)

// DevicesManager 设备管理器（单例）
// 所有 map/slice 在 Init 中一次性构建完成，之后只读，无需锁保护。
// 实时数据的安全由 RtParameter 自持锁保证。
type DevicesManager struct {
	db *gorm.DB

	// 静态配置（Init 后只读）
	devices []*model.DeviceInfo
	units   []*model.UnitInfo
	params  []*model.ParameterInfo

	deviceMap map[string]*model.DeviceInfo
	unitMap   map[string]*model.UnitInfo
	paramMap  map[string]*model.ParameterInfo

	unitsByDevice map[string][]*model.UnitInfo
	paramsByUnit  map[string][]*model.ParameterInfo

	// 实时数据
	rtParams       []*model.RtParameter
	rtParamMap     map[string]*model.RtParameter
	rtParamsByUnit map[string][]*model.RtParameter

	rtDevices   []*model.RtDevice
	rtDeviceMap map[string]*model.RtDevice
}

var (
	instance *DevicesManager
	once     sync.Once
)

// GetDevicesManager 获取单例
func GetDevicesManager() *DevicesManager {
	once.Do(func() {
		instance = &DevicesManager{}
	})
	return instance
}

// Init 初始化：绑定数据库并从数据库加载所有数据到内存（仅启动时调用一次）
func (m *DevicesManager) Init(db *gorm.DB) error {
	m.db = db
	if err := m.reload(); err != nil {
		return err
	}

	// 通过事件总线订阅系统心跳，检测设备在线状态
	GetEventBus().Subscribe(TopicSysTick, m.checkDeviceHeartbeat)
	return nil
}

// reload 从数据库重新加载所有数据到内存
func (m *DevicesManager) reload() error {
	var rawDevices []model.DeviceInfo
	var rawUnits []model.UnitInfo
	var rawParams []model.ParameterInfo

	if err := m.db.Find(&rawDevices).Error; err != nil {
		return err
	}
	if err := m.db.Find(&rawUnits).Error; err != nil {
		return err
	}
	if err := m.db.Find(&rawParams).Error; err != nil {
		return err
	}

	// 构建切片和映射
	m.devices = make([]*model.DeviceInfo, 0, len(rawDevices))
	m.deviceMap = make(map[string]*model.DeviceInfo, len(rawDevices))

	for i := range rawDevices {
		d := &rawDevices[i]
		m.devices = append(m.devices, d)
		m.deviceMap[d.Index] = d
	}

	m.units = make([]*model.UnitInfo, 0, len(rawUnits))
	m.unitMap = make(map[string]*model.UnitInfo, len(rawUnits))
	m.unitsByDevice = make(map[string][]*model.UnitInfo)

	for i := range rawUnits {
		u := &rawUnits[i]
		m.units = append(m.units, u)
		m.unitMap[u.Index] = u
		if idx := strings.IndexByte(u.Index, '.'); idx != -1 {
			parent := u.Index[:idx]
			m.unitsByDevice[parent] = append(m.unitsByDevice[parent], u)
		}
	}

	m.params = make([]*model.ParameterInfo, 0, len(rawParams))
	m.paramMap = make(map[string]*model.ParameterInfo, len(rawParams))
	m.paramsByUnit = make(map[string][]*model.ParameterInfo)

	for i := range rawParams {
		p := &rawParams[i]
		m.params = append(m.params, p)
		m.paramMap[p.Index] = p
		if idx := strings.LastIndexByte(p.Index, '.'); idx != -1 {
			parent := p.Index[:idx]
			m.paramsByUnit[parent] = append(m.paramsByUnit[parent], p)
		}
	}

	// 初始化参数实时数据
	m.rtParams = make([]*model.RtParameter, 0, len(rawParams))
	m.rtParamMap = make(map[string]*model.RtParameter, len(rawParams))
	m.rtParamsByUnit = make(map[string][]*model.RtParameter, len(rawParams))

	for i := range rawParams {
		p := &rawParams[i]
		r := model.NewRtParameter(p.Index, p.Value)
		m.rtParams = append(m.rtParams, r)
		m.rtParamMap[r.Index()] = r
		if idx := strings.LastIndexByte(r.Index(), '.'); idx != -1 {
			parent := r.Index()[:idx]
			m.rtParamsByUnit[parent] = append(m.rtParamsByUnit[parent], r)
		}
	}

	// 初始化设备实时状态
	m.rtDevices = make([]*model.RtDevice, 0, len(rawDevices))
	m.rtDeviceMap = make(map[string]*model.RtDevice, len(rawDevices))

	for i := range rawDevices {
		d := &rawDevices[i]
		r := model.NewRtDevice(d.Index)
		m.rtDevices = append(m.rtDevices, r)
		m.rtDeviceMap[r.DeviceIndex()] = r
	}

	return nil
}

// ==================== 精准访问（O(1)） ====================

// GetDevice 按 Index 获取设备
func (m *DevicesManager) GetDevice(index string) *model.DeviceInfo {
	return m.deviceMap[index]
}

// GetUnit 按 Index 获取单元
func (m *DevicesManager) GetUnit(index string) *model.UnitInfo {
	return m.unitMap[index]
}

// GetParam 按 Index 获取参数配置
func (m *DevicesManager) GetParam(index string) *model.ParameterInfo {
	return m.paramMap[index]
}

// ==================== 遍历访问 ====================

// GetAllDevices 获取所有设备
func (m *DevicesManager) GetAllDevices() []*model.DeviceInfo {
	return m.devices
}

// GetAllUnits 获取所有单元
func (m *DevicesManager) GetAllUnits() []*model.UnitInfo {
	return m.units
}

// GetAllParams 获取所有参数配置
func (m *DevicesManager) GetAllParams() []*model.ParameterInfo {
	return m.params
}

// ==================== 层级访问 ====================

// GetUnitsOfDevice 获取某个设备下的所有单元
func (m *DevicesManager) GetUnitsOfDevice(deviceIndex string) []*model.UnitInfo {
	return m.unitsByDevice[deviceIndex]
}

// GetParamsOfUnit 获取某个单元下的所有参数配置
func (m *DevicesManager) GetParamsOfUnit(unitIndex string) []*model.ParameterInfo {
	return m.paramsByUnit[unitIndex]
}

// ==================== 实时数据访问 ====================
// 返回 RtParameter 指针，读写安全由 RtParameter 自持锁保证。

// GetRtParam 获取单个参数实时数据对象
func (m *DevicesManager) GetRtParam(index string) *model.RtParameter {
	return m.rtParamMap[index]
}

// GetAllRtParams 获取所有参数实时数据对象
func (m *DevicesManager) GetAllRtParams() []*model.RtParameter {
	return m.rtParams
}

// GetRtParamsOfUnit 获取某个单元下所有参数实时数据对象
func (m *DevicesManager) GetRtParamsOfUnit(unitIndex string) []*model.RtParameter {
	return m.rtParamsByUnit[unitIndex]
}

// ==================== 设备实时状态 ====================

// Heartbeat 驱动成功通信时调用，重置设备心跳计数
func (m *DevicesManager) Heartbeat(deviceIndex string) {
	if d := m.rtDeviceMap[deviceIndex]; d != nil {
		d.Heartbeat()
	}
}

// GetRtDevice 获取设备实时状态对象
func (m *DevicesManager) GetRtDevice(deviceIndex string) *model.RtDevice {
	return m.rtDeviceMap[deviceIndex]
}

// checkDeviceHeartbeat 系统定时器回调：检测所有设备心跳超时
func (m *DevicesManager) checkDeviceHeartbeat() {
	const maxCount = 24 // 24 秒无心跳 → 离线，可调整
	for _, d := range m.rtDevices {
		if changed := d.TickAndCheck(maxCount); changed {
			// 状态变更，后续可在此触发告警或日志
		}
	}
}
