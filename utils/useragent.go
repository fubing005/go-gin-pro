package utils

import (
	"github.com/mssola/user_agent"
)

// UserAgentInfo 用户代理信息
type UserAgentInfo struct {
	Device  string // 设备类型
	Browser string // 浏览器
	OS      string // 操作系统
}

// ParseUserAgent 解析User-Agent
func ParseUserAgent(userAgent string) *UserAgentInfo {
	ua := user_agent.New(userAgent)

	// 获取设备类型
	device := "Desktop"
	if ua.Mobile() {
		device = "Mobile"
	}

	// 获取浏览器信息
	browser, version := ua.Browser()

	// 获取操作系统信息
	os := ua.OS()

	return &UserAgentInfo{
		Device:  device,
		Browser: browser + " " + version,
		OS:      os,
	}
}

// GetLocationByIP 根据IP获取地理位置
func GetLocationByIP(ip string) string {
	// TODO: 接入IP地址库或第三方服务获取地理位置
	// 这里先返回空字符串
	return "需要接入第三方服务"
}
