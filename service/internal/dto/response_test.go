package dto

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Success(c, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态码 = %d, want 200", w.Code)
	}

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("Message = %s, want success", resp.Message)
	}
}

func TestSuccessPage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	data := []string{"a", "b"}
	SuccessPage(c, data, 1, 20, 100)

	var resp PageResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Meta.Page != 1 {
		t.Errorf("Meta.Page = %d, want 1", resp.Meta.Page)
	}
	if resp.Meta.PageSize != 20 {
		t.Errorf("Meta.PageSize = %d, want 20", resp.Meta.PageSize)
	}
	if resp.Meta.Total != 100 {
		t.Errorf("Meta.Total = %d, want 100", resp.Meta.Total)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Error(c, ErrCodeParamInvalid, "参数错误")

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != ErrCodeParamInvalid {
		t.Errorf("Code = %d, want %d", resp.Code, ErrCodeParamInvalid)
	}
	if resp.Message != "参数错误" {
		t.Errorf("Message = %s, want 参数错误", resp.Message)
	}
}

func TestErrorWithStatus(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	ErrorWithStatus(c, http.StatusForbidden, ErrCodeForbidden, "无权限")

	if w.Code != http.StatusForbidden {
		t.Errorf("HTTP 状态码 = %d, want 403", w.Code)
	}
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{ErrCodeSuccess, "操作成功"},
		{ErrCodeUnauthorized, "未登录或登录已过期"},
		{ErrCodeEntityNotFound, "主体不存在"},
		{ErrCodeVideoNotFound, "视频不存在"},
		{99999, "未知错误"},
	}

	for _, tt := range tests {
		got := GetMessage(tt.code)
		if got != tt.expected {
			t.Errorf("GetMessage(%d) = %s, want %s", tt.code, got, tt.expected)
		}
	}
}

func TestPagination_Offset(t *testing.T) {
	p := Pagination{Page: 3, PageSize: 20}
	if p.Offset() != 40 {
		t.Errorf("Offset() = %d, want 40", p.Offset())
	}

	p2 := Pagination{Page: 1, PageSize: 10}
	if p2.Offset() != 0 {
		t.Errorf("Offset() = %d, want 0", p2.Offset())
	}
}

func TestDefaultPagination(t *testing.T) {
	p := DefaultPagination()
	if p.Page != 1 {
		t.Errorf("Page = %d, want 1", p.Page)
	}
	if p.PageSize != 20 {
		t.Errorf("PageSize = %d, want 20", p.PageSize)
	}
}
