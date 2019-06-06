package helpers

import (
	"net"
	"net/http"
	"strings"
)

// GetRemoteIP 返回远程客户端IP
func GetRemoteIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		// 多层代理情况的处理
		remoteAddr = strings.Split(ip, ",")[0]
	} else if ip = req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
