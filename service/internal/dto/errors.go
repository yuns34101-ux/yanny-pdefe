package dto

// 错误码定义

const (
	// 通用错误 10001-19999
	ErrCodeSuccess      = 0
	ErrCodeParamInvalid = 10001
	ErrCodeUnauthorized = 10002
	ErrCodeForbidden    = 10003
	ErrCodeNotFound     = 10004
	ErrCodeInternal     = 10005
	ErrCodeDuplicate    = 10006

	// 主体/小程序 20001-29999
	ErrCodeEntityNotFound   = 20001
	ErrCodeMpAccountNotFound = 20002

	// 视频/分类 30001-39999
	ErrCodeVideoNotFound     = 30001
	ErrCodeCategoryNotFound  = 30002
	ErrCodeCommentNotFound   = 30003
	ErrCodeVideoOffline      = 30004

	// 用户/管理后台 40001-49999
	ErrCodeUserNotFound     = 40001
	ErrCodeAdminNotFound    = 40002
	ErrCodePasswordWrong    = 40003
	ErrCodeAdminDisabled    = 40004
	ErrCodeLoginFailed      = 40005
	ErrCodeRoleNotFound     = 40006
	ErrCodePermissionDenied = 40007

	// CDN/文件 50001-59999
	ErrCodeUploadFailed = 50001
)

// 错误消息映射
var errMessages = map[int]string{
	ErrCodeSuccess:            "操作成功",
	ErrCodeParamInvalid:       "参数无效",
	ErrCodeUnauthorized:       "未登录或登录已过期",
	ErrCodeForbidden:          "无权限访问",
	ErrCodeNotFound:           "资源不存在",
	ErrCodeInternal:           "服务器内部错误",
	ErrCodeDuplicate:          "数据已存在",
	ErrCodeEntityNotFound:     "主体不存在",
	ErrCodeMpAccountNotFound:  "小程序账号不存在",
	ErrCodeVideoNotFound:      "视频不存在",
	ErrCodeVideoOffline:       "视频已下架",
	ErrCodeCategoryNotFound:   "分类不存在",
	ErrCodeCommentNotFound:    "评论不存在",
	ErrCodeUserNotFound:       "用户不存在",
	ErrCodeAdminNotFound:      "管理员不存在",
	ErrCodePasswordWrong:      "密码错误",
	ErrCodeAdminDisabled:      "管理员账号已被禁用",
	ErrCodeLoginFailed:        "登录失败",
	ErrCodeRoleNotFound:       "角色不存在",
	ErrCodePermissionDenied:   "权限不足",
	ErrCodeUploadFailed:       "文件上传失败",
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	if msg, ok := errMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
