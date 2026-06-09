package model

import (
	"sync"
	"time"
)

type SCADAObject struct {
	Index       string    `gorm:"column:index;type:varchar(256);primaryKey;comment:索引,全局唯一层级命名" json:"index"`
	Name        string    `gorm:"column:name;type:varchar(128);not null;comment:名称" json:"name"`
	Description string    `gorm:"column:description;type:varchar(512);default:'';comment:描述" json:"description"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"update_time"`
}

// Task 采集任务（一个任务 = 一个采集协程，关联一个设备）
type Task struct {
	SCADAObject
	DeviceIndex string `gorm:"column:device_index;type:varchar(256);index;not null;comment:关联设备Index" json:"device_index"`
	Driver      string `gorm:"column:driver;type:varchar(32);not null;comment:驱动名称: modbus,s7,iec104" json:"driver"`
	Transport   string `gorm:"column:transport;type:varchar(16);not null;comment:传输方式: tcp,serial" json:"transport"`
	Config      string `gorm:"column:config;type:text;default:'{}';comment:采集配置(JSON),驱动自行解析" json:"config"`
	Enable      bool   `gorm:"column:enable;not null;default:true;comment:是否启用" json:"enable"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	SCADAObject
	Model        string    `gorm:"column:model;type:varchar(64);default:'';comment:设备型号" json:"model"`
	Manufacturer string    `gorm:"column:manufacturer;type:varchar(128);default:'';comment:生产厂家" json:"manufacturer"`
	SerialNumber string    `gorm:"column:serial_number;type:varchar(64);default:'';comment:序列号" json:"serial_number"`
	Location     string    `gorm:"column:location;type:varchar(256);default:'';comment:安装位置" json:"location"`
	Owner        string    `gorm:"column:owner;type:varchar(64);default:'';comment:所属用户" json:"owner"`
	Permission   string    `gorm:"column:permission;type:varchar(64);default:'';comment:权限: r-只读 rw-读写" json:"permission"`
	InstallDate  time.Time `gorm:"column:install_date;comment:安装日期" json:"install_date"`
}

// UnitInfo 单元信息
type UnitInfo struct {
	SCADAObject
}

// ParameterInfo 参数信息
type ParameterInfo struct {
	SCADAObject
	DataType   string  `gorm:"column:data_type;type:varchar(16);default:'';comment:数据类型" json:"data_type"`
	Value      string  `gorm:"column:value;type:varchar(64);default:'';comment:默认值" json:"value"`
	Unit       string  `gorm:"column:unit;type:varchar(16);default:'';comment:单位" json:"unit"`
	MinValue   float64 `gorm:"column:min_value;type:double;default:0;comment:量程下限" json:"min_value"`
	MaxValue   float64 `gorm:"column:max_value;type:double;default:0;comment:量程上限" json:"max_value"`
	AlarmUpper float64 `gorm:"column:alarm_upper;type:double;default:0;comment:报警上限" json:"alarm_upper"`
	AlarmLower float64 `gorm:"column:alarm_lower;type:double;default:0;comment:报警下限" json:"alarm_lower"`
	Type       int     `gorm:"column:type;not null;default:0;comment:参数类型(0-采集 1-计算 2-常量)" json:"type"`
	UploadType int     `gorm:"column:upload_type;not null;default:0;comment:上传类型(0-变化上报 1-上升沿 2-下降沿 3-定时)" json:"upload_type"`
	UploadInterval int `gorm:"column:upload_interval;not null;default:0;comment:上报间隔(秒)" json:"upload_interval"`
}

// RtDevice 设备实时状态（自持锁，安全并发读写）
type RtDevice struct {
	mu               sync.Mutex
	deviceIndex      string
	online           bool
	connectTime      int64 // 连接时间（毫秒时间戳）
	heartbeatCounter int   // 心跳计数，驱动成功轮询时清零，定时器每秒+1
}

// NewRtDevice 创建设备实时状态
func NewRtDevice(deviceIndex string) *RtDevice {
	return &RtDevice{
		deviceIndex: deviceIndex,
	}
}

// DeviceIndex 返回设备索引
func (d *RtDevice) DeviceIndex() string {
	return d.deviceIndex
}

// Snapshot 返回只读副本
func (d *RtDevice) Snapshot() RtDeviceSnapshot {
	d.mu.Lock()
	defer d.mu.Unlock()
	return RtDeviceSnapshot{
		DeviceIndex: d.deviceIndex,
		Online:      d.online,
		ConnectTime: d.connectTime,
	}
}

// Heartbeat 心跳：驱动成功通信时调用，清零计数器
func (d *RtDevice) Heartbeat() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.online {
		d.online = true
		d.connectTime = time.Now().UnixMilli()
	}
	d.heartbeatCounter = 0
}

// TickAndCheck 定时器回调：计数器+1，超时判定离线
func (d *RtDevice) TickAndCheck(maxCount int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.heartbeatCounter++
	if d.online && d.heartbeatCounter > maxCount {
		d.online = false
		return true // 状态变为离线
	}
	return false
}

// RtDeviceSnapshot 设备实时状态快照
type RtDeviceSnapshot struct {
	DeviceIndex string `json:"device_index"`
	Online      bool   `json:"online"`
	ConnectTime int64  `json:"connect_time"`
}

// RtSnapshot 实时参数快照（只读副本）
type RtSnapshot struct {
	Index     string `json:"index"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
	Quality   int    `json:"quality"`
}

// RtParameter 实时参数（字段私有，自持锁，安全并发读写）
type RtParameter struct {
	mu        sync.Mutex
	index     string
	value     string
	timestamp int64
	quality   int
}

// NewRtParameter 创建实时参数实例
func NewRtParameter(index, initialValue string) *RtParameter {
	return &RtParameter{
		index: index,
		value: initialValue,
	}
}

// Index 返回索引（创建后不变，无锁安全）
func (r *RtParameter) Index() string {
	return r.index
}

// Get 返回数据副本（线程安全）
func (r *RtParameter) Get() RtSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return RtSnapshot{
		Index:     r.index,
		Value:     r.value,
		Timestamp: r.timestamp,
		Quality:   r.quality,
	}
}

// Update 安全更新实时值
func (r *RtParameter) Update(value string, quality int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.value = value
	r.quality = quality
	r.timestamp = time.Now().UnixMilli()
}