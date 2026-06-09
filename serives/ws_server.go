package serives

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"SH_SCADA/common/model"
	"SH_SCADA/plat"

	"github.com/gorilla/websocket"
)

// ==================== 消息协议：type 字段路由 ====================

// Message 通用消息帧，type 决定 data 的结构
type Message struct {
	Type      string          `json:"type"`                 // 消息类型
	RequestID string          `json:"request_id,omitempty"` // 请求 ID
	Data      json.RawMessage `json:"data,omitempty"`       // 数据载荷
}

// ==================== type 常量定义 ====================

// 客户端 → 服务端
const (
	TypeSubscribe       = "subscribe"
	TypeUnsubscribe     = "unsubscribe"
	TypeCommand         = "command"
	TypeGetDevices      = "get_devices"
	TypeGetUnits        = "get_units"
	TypeGetParams       = "get_params"
	TypeGetDeviceStatus = "get_device_status"
)

// 服务端 → 客户端
const (
	TypePushRtData   = "push_rtdata"
	TypeDeviceList   = "device_list"
	TypeUnitList     = "unit_list"
	TypeParamList    = "param_list"
	TypeDeviceStatus = "device_status"
	TypeResponse     = "response"
)

// ==================== 各 type 对应的 data 结构体 ====================

type SubscribeData struct {
	Device string `json:"device"`
}

type CommandData struct {
	Index string `json:"index"`
	Value string `json:"value"`
}

type GetUnitsData struct {
	Device string `json:"device"`
}

type GetParamsData struct {
	Unit string `json:"unit"`
}

type GetDeviceStatusData struct {
	Device string `json:"device"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ==================== Hub 管理所有连接 ====================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHub WebSocket 连接管理中心
type WSHub struct {
	clients    map[*wsClient]bool
	register   chan *wsClient
	unregister chan *wsClient
	mu         sync.RWMutex
}

// wsClient 单个 WebSocket 连接
type wsClient struct {
	conn       *websocket.Conn
	send       chan []byte
	hub        *WSHub
	subscribed map[string]bool
	mu         sync.Mutex
}

var (
	wsHubInstance *WSHub
	wsHubOnce     sync.Once
)

// GetWSHub 获取单例
func GetWSHub() *WSHub {
	wsHubOnce.Do(func() {
		wsHubInstance = &WSHub{
			clients:    make(map[*wsClient]bool),
			register:   make(chan *wsClient),
			unregister: make(chan *wsClient),
		}
	})
	return wsHubInstance
}

// HandleWS WebSocket 升级处理（token 校验通过后才允许连接）
func (h *WSHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	claims, err := ValidateWSToken(r)
	if err != nil {
		http.Error(w, "token 无效或已过期", 401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] 升级失败 (%s): %v", claims.Username, err)
		return
	}

	log.Printf("[WS] 客户端连接: %s (%s)", claims.Username, claims.Role)
	c := &wsClient{
		conn:       conn,
		send:       make(chan []byte, 64),
		hub:        h,
		subscribed: make(map[string]bool),
	}
	h.register <- c
	go c.writePump()
	go c.readPump()
}

// Run 运行 Hub（需在 goroutine 中调用）
func (h *WSHub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()
		}
	}
}

// PushData 定时推送实时数据（TopicSysTick 事件订阅）
func (h *WSHub) PushData() {
	mgr := plat.GetDevicesManager()

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		c.mu.Lock()
		devices := make([]string, 0, len(c.subscribed))
		for dev := range c.subscribed {
			devices = append(devices, dev)
		}
		c.mu.Unlock()

		if len(devices) == 0 {
			continue
		}

		var snapshots []model.RtSnapshot
		for _, dev := range devices {
			for _, p := range mgr.GetRtParamsOfUnit(dev) {
				snapshots = append(snapshots, p.Get())
			}
		}

		data, _ := json.Marshal(snapshots)
		c.sendMsg(TypePushRtData, "", data)
	}
}

// ==================== 消息分发 ====================

func (c *wsClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			c.replyError(msg.RequestID, "消息格式错误")
			continue
		}

		switch msg.Type {
		case TypeSubscribe:
			c.handleSubscribe(msg)
		case TypeUnsubscribe:
			c.handleUnsubscribe(msg)
		case TypeCommand:
			c.handleCommand(msg)
		case TypeGetDevices:
			c.handleGetDevices(msg)
		case TypeGetUnits:
			c.handleGetUnits(msg)
		case TypeGetParams:
			c.handleGetParams(msg)
		case TypeGetDeviceStatus:
			c.handleGetDeviceStatus(msg)
		default:
			c.replyError(msg.RequestID, fmt.Sprintf("未知消息类型: %s", msg.Type))
		}
	}
}

func (c *wsClient) writePump() {
	defer c.conn.Close()
	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}

// ==================== 消息处理器 ====================

func (c *wsClient) handleSubscribe(msg Message) {
	var d SubscribeData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	c.mu.Lock()
	c.subscribed[d.Device] = true
	c.mu.Unlock()
	c.replyOk(msg.RequestID)
}

func (c *wsClient) handleUnsubscribe(msg Message) {
	var d SubscribeData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	c.mu.Lock()
	delete(c.subscribed, d.Device)
	c.mu.Unlock()
	c.replyOk(msg.RequestID)
}

func (c *wsClient) handleCommand(msg Message) {
	var d CommandData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	fmt.Printf("[WS] 收到命令: index=%s value=%s\n", d.Index, d.Value)
	plat.GetEventBus().Publish(plat.TopicCommand, d.Index, d.Value)
	c.replyOk(msg.RequestID)
}

func (c *wsClient) handleGetDevices(msg Message) {
	devices := plat.GetDevicesManager().GetAllDevices()
	data, _ := json.Marshal(devices)
	c.sendMsg(TypeDeviceList, msg.RequestID, data)
}

func (c *wsClient) handleGetUnits(msg Message) {
	var d GetUnitsData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	units := plat.GetDevicesManager().GetUnitsOfDevice(d.Device)
	data, _ := json.Marshal(units)
	c.sendMsg(TypeUnitList, msg.RequestID, data)
}

func (c *wsClient) handleGetParams(msg Message) {
	var d GetParamsData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	params := plat.GetDevicesManager().GetParamsOfUnit(d.Unit)
	data, _ := json.Marshal(params)
	c.sendMsg(TypeParamList, msg.RequestID, data)
}

func (c *wsClient) handleGetDeviceStatus(msg Message) {
	var d GetDeviceStatusData
	if err := json.Unmarshal(msg.Data, &d); err != nil {
		c.replyError(msg.RequestID, "参数错误")
		return
	}
	rt := plat.GetDevicesManager().GetRtDevice(d.Device)
	if rt == nil {
		c.replyError(msg.RequestID, "设备不存在")
		return
	}
	snap := rt.Snapshot()
	data, _ := json.Marshal(snap)
	c.sendMsg(TypeDeviceStatus, msg.RequestID, data)
}

// ==================== 发送工具方法 ====================

func (c *wsClient) sendMsg(msgType, requestID string, data json.RawMessage) {
	msg := Message{
		Type:      msgType,
		RequestID: requestID,
		Data:      data,
	}
	raw, _ := json.Marshal(msg)
	select {
	case c.send <- raw:
	default:
	}
}

func (c *wsClient) replyOk(requestID string) {
	resp, _ := json.Marshal(ResponseData{Success: true})
	c.sendMsg(TypeResponse, requestID, resp)
}

func (c *wsClient) replyError(requestID, errMsg string) {
	resp, _ := json.Marshal(ResponseData{Success: false, Error: errMsg})
	c.sendMsg(TypeResponse, requestID, resp)
}
