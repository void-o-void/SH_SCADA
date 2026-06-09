package serives

import (
	"encoding/json"
	"net/http"
	"strings"

	"SH_SCADA/common/db"
	"SH_SCADA/common/model"
	"SH_SCADA/plat"
)

// RegisterREST 注册 REST API 路由
func RegisterREST() {
	// 登录（无需鉴权）
	http.HandleFunc("/api/login", HandleLogin)

	// 设备 CRUD（需鉴权）
	http.HandleFunc("/api/devices", AuthMiddleware(handleDevices))
	http.HandleFunc("/api/devices/", AuthMiddleware(handleDeviceByIndex))

	// 单元 CRUD
	http.HandleFunc("/api/units", AuthMiddleware(handleUnits))
	http.HandleFunc("/api/units/", AuthMiddleware(handleUnitByIndex))

	// 参数 CRUD
	http.HandleFunc("/api/params", AuthMiddleware(handleParams))
	http.HandleFunc("/api/params/", AuthMiddleware(handleParamByIndex))

	// 任务 CRUD
	http.HandleFunc("/api/tasks", AuthMiddleware(handleTasks))
	http.HandleFunc("/api/tasks/", AuthMiddleware(handleTaskByIndex))
}

// ==================== 设备 ====================

func handleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list(plat.GetDevicesManager().GetAllDevices(), w)
	case http.MethodPost:
		crudSave[model.DeviceInfo](r, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func handleDeviceByIndex(w http.ResponseWriter, r *http.Request) {
	index := extractIndex(r.URL.Path, "/api/devices/")
	switch r.Method {
	case http.MethodGet:
		d := plat.GetDevicesManager().GetDevice(index)
		if d == nil {
			http.NotFound(w, r)
			return
		}
		jsonOK(w, d)
	case http.MethodPut:
		crudUpdate[model.DeviceInfo](index, r, w)
	case http.MethodDelete:
		crudDelete[model.DeviceInfo](index, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// ==================== 单元 ====================

func handleUnits(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list(plat.GetDevicesManager().GetAllUnits(), w)
	case http.MethodPost:
		crudSave[model.UnitInfo](r, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func handleUnitByIndex(w http.ResponseWriter, r *http.Request) {
	index := extractIndex(r.URL.Path, "/api/units/")
	switch r.Method {
	case http.MethodGet:
		u := plat.GetDevicesManager().GetUnit(index)
		if u == nil {
			http.NotFound(w, r)
			return
		}
		jsonOK(w, u)
	case http.MethodPut:
		crudUpdate[model.UnitInfo](index, r, w)
	case http.MethodDelete:
		crudDelete[model.UnitInfo](index, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// ==================== 参数 ====================

func handleParams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list(plat.GetDevicesManager().GetAllParams(), w)
	case http.MethodPost:
		crudSave[model.ParameterInfo](r, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func handleParamByIndex(w http.ResponseWriter, r *http.Request) {
	index := extractIndex(r.URL.Path, "/api/params/")
	switch r.Method {
	case http.MethodGet:
		p := plat.GetDevicesManager().GetParam(index)
		if p == nil {
			http.NotFound(w, r)
			return
		}
		jsonOK(w, p)
	case http.MethodPut:
		crudUpdate[model.ParameterInfo](index, r, w)
	case http.MethodDelete:
		crudDelete[model.ParameterInfo](index, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// ==================== 任务 ====================

func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list(plat.GetTaskManager().GetAllTasks(), w)
	case http.MethodPost:
		crudSave[model.Task](r, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func handleTaskByIndex(w http.ResponseWriter, r *http.Request) {
	index := extractIndex(r.URL.Path, "/api/tasks/")
	switch r.Method {
	case http.MethodGet:
		t := plat.GetTaskManager().GetTask(index)
		if t == nil {
			http.NotFound(w, r)
			return
		}
		jsonOK(w, t)
	case http.MethodPut:
		crudUpdate[model.Task](index, r, w)
	case http.MethodDelete:
		crudDelete[model.Task](index, w)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// ==================== 通用 CRUD 工具 ====================

func crudSave[T any](r *http.Request, w http.ResponseWriter) {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "JSON 解析失败", 400)
		return
	}
	if err := db.DB.Save(&obj).Error; err != nil {
		http.Error(w, "保存失败", 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func crudUpdate[T any](index string, r *http.Request, w http.ResponseWriter) {
	var obj T
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "JSON 解析失败", 400)
		return
	}
	if err := db.DB.Where("index = ?", index).Updates(&obj).Error; err != nil {
		http.Error(w, "更新失败", 500)
		return
	}
	jsonOK(w, map[string]string{"status": "ok"})
}

func crudDelete[T any](index string, w http.ResponseWriter) {
	if err := db.DB.Where("index = ?", index).Delete(new(T)).Error; err != nil {
		http.Error(w, "删除失败", 500)
		return
	}
	jsonOK(w, map[string]string{"status": "ok"})
}

func list[T any](data []*T, w http.ResponseWriter) {
	if data == nil {
		data = []*T{}
	}
	jsonOK(w, data)
}

// ==================== 工具 ====================

func extractIndex(path, prefix string) string {
	index := strings.TrimPrefix(path, prefix)
	return strings.TrimSuffix(index, "/")
}

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
