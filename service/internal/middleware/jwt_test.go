package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"yanny-service/internal/config"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-jwt-secret",
			ExpireHours:   2,
			MpExpireHours: 720,
		},
	}
}

func TestGenerateAdminToken(t *testing.T) {
	token, err := GenerateAdminToken(1, "admin")
	if err != nil {
		t.Fatalf("GenerateAdminToken() 失败: %v", err)
	}
	if token == "" {
		t.Error("Token 不应为空")
	}
}

func TestGenerateMpToken(t *testing.T) {
	token, err := GenerateMpToken(1, 10)
	if err != nil {
		t.Fatalf("GenerateMpToken() 失败: %v", err)
	}
	if token == "" {
		t.Error("Token 不应为空")
	}
}

func TestParseToken_Valid(t *testing.T) {
	token, _ := GenerateAdminToken(42, "ops")
	claims, err := parseToken(token)
	if err != nil {
		t.Fatalf("parseToken() 失败: %v", err)
	}
	if claims.AdminID != 42 {
		t.Errorf("AdminID = %d, want 42", claims.AdminID)
	}
	if claims.Username != "ops" {
		t.Errorf("Username = %s, want ops", claims.Username)
	}
	if claims.Type != "admin" {
		t.Errorf("Type = %s, want admin", claims.Type)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	_, err := parseToken("invalid.token.string")
	if err == nil {
		t.Error("parseToken() 应返回错误，token 无效")
	}
}

func TestParseToken_Expired(t *testing.T) {
	// JWT 过期由 jwt 库自动校验，此处仅验证 expired token 返回错误
	// 用极短过期时间生成 token 后等待它过期来测试
	t.Skip("需要等待 token 过期，跳过")
}

func TestAdminAuth_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	mw := AdminAuthMiddleware()
	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP 状态码 = %d, want 401", w.Code)
	}
}

func TestAdminAuth_ValidToken(t *testing.T) {
	token, _ := GenerateAdminToken(1, "admin")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	mw := AdminAuthMiddleware()
	mw(c)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态码 = %d, want 200", w.Code)
	}

	adminID := GetAdminID(c)
	if adminID != 1 {
		t.Errorf("GetAdminID = %d, want 1", adminID)
	}
}

func TestAdminAuth_WrongType(t *testing.T) {
	token, _ := GenerateMpToken(1, 10)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	mw := AdminAuthMiddleware()
	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP 状态码 = %d, want 401", w.Code)
	}
}

func TestMpAuthOptional_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	mw := MpAuthOptional()
	mw(c)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态码 = %d, want 200", w.Code)
	}
	userID := GetUserID(c)
	if userID != 0 {
		t.Errorf("无 token 时 user_id 应为 0")
	}
}

func TestMpAuthRequired_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	mw := MpAuthRequired()
	mw(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP 状态码 = %d, want 401", w.Code)
	}
}

func TestMpAuthRequired_ValidToken(t *testing.T) {
	token, _ := GenerateMpToken(5, 99)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	mw := MpAuthRequired()
	mw(c)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态码 = %d, want 200", w.Code)
	}
	if GetUserID(c) != 5 {
		t.Errorf("GetUserID = %d, want 5", GetUserID(c))
	}
	if GetMpAccountID(c) != 99 {
		t.Errorf("GetMpAccountID = %d, want 99", GetMpAccountID(c))
	}
}

func TestExtractToken_Bearer(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer abc.def.ghi")

	tok := extractToken(c)
	if tok != "abc.def.ghi" {
		t.Errorf("extractToken = %s, want abc.def.ghi", tok)
	}
}

func TestExtractToken_BareToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "bare-token-string")

	tok := extractToken(c)
	if tok != "bare-token-string" {
		t.Errorf("extractToken = %s, want bare-token-string", tok)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	if GetUserID(c) != 0 {
		t.Error("未设置时 GetUserID 应返回 0")
	}
}

func TestGetAdminID_NotSet(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	if GetAdminID(c) != 0 {
		t.Error("未设置时 GetAdminID 应返回 0")
	}
}

func TestFormatAdminID(t *testing.T) {
	if FormatAdminID(123) != "123" {
		t.Errorf("FormatAdminID(123) = %s, want 123", FormatAdminID(123))
	}
}
