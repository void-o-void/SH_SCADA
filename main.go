package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"SH_SCADA/common/db"
	"SH_SCADA/plat"
	"SH_SCADA/serives"

	// 按需引入驱动
	// "SH_SCADA/drivers/modbus"
	// "SH_SCADA/drivers/iec104"
)

func main() {
	// 1. 初始化数据库
	if err := db.Init("scada.db"); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 2. 加载静态配置到内存
	if err := plat.GetDevicesManager().Init(db.DB); err != nil {
		log.Fatalf("设备管理器初始化失败: %v", err)
	}

	// 3. 注册采集驱动（按需引入）
	tm := plat.GetTaskManager()
	// tm.RegisterDriver(modbus.New())
	// tm.RegisterDriver(iec104.New())

	// 4. 加载任务并启动采集
	if err := tm.Init(db.DB); err != nil {
		log.Fatalf("任务管理器初始化失败: %v", err)
	}
	tm.StartAll(plat.GetDevicesManager())

	// 5. 启动 WebSocket Hub
	hub := serives.GetWSHub()
	go hub.Run()

	// 6. 通过事件总线订阅：WS 实时数据推送、设备心跳检测
	//    （DevicesManager.Init 中已订阅 TopicSysTick，此处订阅 WS 推送）
	plat.GetEventBus().Subscribe(plat.TopicSysTick, hub.PushData)

	// 7. 启动系统定时器（通过事件总线发布 TopicSysTick）
	timer := plat.GetSysTimer()
	timer.Start()

	// 8. HTTP 路由：REST API + WebSocket
	serives.RegisterREST()
	http.HandleFunc("/ws", hub.HandleWS)
	go func() {
		addr := "0.0.0.0:8080"
		fmt.Printf("WebSocket 服务启动: ws://%s/ws\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	fmt.Println("SH_SCADA 网关启动成功")

	// 9. 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭...")
	timer.Stop()
	tm.StopAll()
	fmt.Println("SH_SCADA 网关已退出")
}