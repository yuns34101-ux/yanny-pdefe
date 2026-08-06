package util

import (
	"testing"
	"time"
)

func TestGenerateDynamicCode(t *testing.T) {
	t1 := time.Date(2026, 8, 6, 14, 30, 0, 0, time.UTC)
	code := GenerateDynamicCode("", t1)
	if code != "1430" {
		t.Errorf("GenerateDynamicCode = %s, want 1430", code)
	}

	t2 := time.Date(2026, 8, 6, 9, 5, 0, 0, time.UTC)
	code = GenerateDynamicCode("", t2)
	if code != "0905" {
		t.Errorf("GenerateDynamicCode = %s, want 0905", code)
	}
}

func TestValidateDynamicCode(t *testing.T) {
	now := time.Now()
	code := GenerateDynamicCode("", now)
	if !ValidateDynamicCode("", code, now) {
		t.Error("当前分钟应通过")
	}

	prevMin := GenerateDynamicCode("", now.Add(-1*time.Minute))
	if !ValidateDynamicCode("", prevMin, now) {
		t.Error("前一分钟应通过")
	}

	nextMin := GenerateDynamicCode("", now.Add(1*time.Minute))
	if !ValidateDynamicCode("", nextMin, now) {
		t.Error("后一分钟应通过")
	}

	old := GenerateDynamicCode("", now.Add(-2*time.Minute))
	if ValidateDynamicCode("", old, now) {
		t.Error("两分钟前不应通过")
	}
}

func TestValidateDynamicCode_Empty(t *testing.T) {
	if ValidateDynamicCode("", "", time.Now()) {
		t.Error("空口令应返回 false")
	}
}

func TestSplitPasswordAndCode(t *testing.T) {
	tests := []struct {
		input        string
		wantPassword string
		wantCode     string
	}{
		{"mypassword1430", "mypassword", "1430"},
		{"abc0905", "abc", "0905"},
		{"1234", "", "1234"},
	}

	for _, tt := range tests {
		password, code := SplitPasswordAndCode(tt.input)
		if password != tt.wantPassword || code != tt.wantCode {
			t.Errorf("SplitPasswordAndCode(%q) = (%q, %q), want (%q, %q)",
				tt.input, password, code, tt.wantPassword, tt.wantCode)
		}
	}
}

func TestGetCurrentTimeHint(t *testing.T) {
	hint := GetCurrentTimeHint()
	if hint == "" {
		t.Error("时间提示不应为空")
	}
}
