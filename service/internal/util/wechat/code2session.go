package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// Code2SessionResult 微信 code2Session 接口返回结果
type Code2SessionResult struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 用 wx.login 获取的 code 换取 openid/session_key
func Code2Session(appid, secret, code string) (*Code2SessionResult, error) {
	query := url.Values{}
	query.Set("appid", appid)
	query.Set("secret", secret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(code2SessionURL + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Code2SessionResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信 code2Session 失败: errcode=%d errmsg=%s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}
