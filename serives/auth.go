package serives

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecret    = "scada-gateway-secret-change-in-production"
	jwtExpire    = 24 * time.Hour
	TokenQuery   = "token"
	AuthHeader   = "Authorization"
	BearerPrefix = "Bearer "
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Claims JWT 载荷
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 内置用户（生产环境应存数据库和哈希密码）
var builtinUsers = map[string]struct {
	passwordHash string
	role         string
}{
	"admin":    {hashPassword("admin123"), "admin"},
	"operator": {hashPassword("oper123"), "operator"},
	"viewer":   {hashPassword("view123"), "viewer"},
}

// HandleLogin 登录接口
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", 400)
		return
	}

	user, ok := builtinUsers[req.Username]
	if !ok || user.passwordHash != hashPassword(req.Password) {
		http.Error(w, "用户名或密码错误", 403)
		return
	}

	token, err := generateToken(req.Username, user.role)
	if err != nil {
		http.Error(w, "生成 token 失败", 500)
		return
	}

	jsonOK(w, LoginResponse{
		Token:    token,
		Username: req.Username,
		Role:     user.role,
	})
}

// AuthMiddleware 为 REST API 加鉴权
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "缺少认证 token", 401)
			return
		}
		if _, err := validateToken(token); err != nil {
			http.Error(w, "token 无效或已过期", 401)
			return
		}
		next(w, r)
	}
}

// ValidateWSToken 校验 WebSocket 连接中的 token
func ValidateWSToken(r *http.Request) (*Claims, error) {
	token := r.URL.Query().Get(TokenQuery)
	if token == "" {
		return nil, jwt.ErrSignatureInvalid
	}
	return validateToken(token)
}

// ==================== 内部工具 ====================

func generateToken(username, role string) (string, error) {
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(jwtSecret))
}

func validateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func extractToken(r *http.Request) string {
	// 优先从 Header 取
	auth := r.Header.Get(AuthHeader)
	if strings.HasPrefix(auth, BearerPrefix) {
		return strings.TrimPrefix(auth, BearerPrefix)
	}
	// 也支持 query 参数（WebSocket 走这个）
	return r.URL.Query().Get(TokenQuery)
}

func hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(hash[:])
}
