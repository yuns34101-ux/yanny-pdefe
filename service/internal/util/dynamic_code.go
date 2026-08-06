package util

import (
	"fmt"
	"time"
)

// DynamicCode 动态口令工具
// 简化方案：口令 = 当前时分（HHMM），普通运营人员看时钟即可获知
// 有效期：当前分钟，允许 ±1 分钟偏差
// 配合防爆破（5 次失败锁 30 分钟）保障安全

const dynamicCodeLength = 4

// GenerateDynamicCode 生成当前分钟的动态口令（纯时钟 HHMM）
// key 参数保留兼容性，当前不使用
func GenerateDynamicCode(key string, t time.Time) string {
	return t.Format("1504") // HHMM
}

// ValidateDynamicCode 校验动态口令是否有效
// 允许当前分钟 ±1 分钟偏差（共 3 个窗口）
// key 参数保留兼容性，当前不使用
func ValidateDynamicCode(key, code string, t time.Time) bool {
	if code == "" {
		return false
	}

	for offset := -1; offset <= 1; offset++ {
		candidate := GenerateDynamicCode(key, t.Add(time.Duration(offset)*time.Minute))
		if candidate == code {
			return true
		}
	}
	return false
}

// GetCurrentTimeHint 获取当前时间提示（登录页展示用）
func GetCurrentTimeHint() string {
	now := time.Now()
	return fmt.Sprintf("当前服务器时间 %s，动态口令为 %s",
		now.Format("15:04"),
		now.Format("1504"))
}

// SplitPasswordAndCode 将用户输入分离为静态密码和动态口令
// 约定：最后 4 位为动态口令（HHMM），前面部分为静态密码
func SplitPasswordAndCode(input string) (password, code string) {
	if len(input) <= dynamicCodeLength {
		return "", input
	}
	return input[:len(input)-dynamicCodeLength], input[len(input)-dynamicCodeLength:]
}
