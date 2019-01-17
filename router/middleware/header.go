package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	// 设置为*的时候前端不会发送cookie过来
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	// 禁用了客户端的 MIME 类型嗅探行为
	c.Header("X-Content-Type-Options", "nosniff")
	// 启用XSS过滤,如果检测到攻击，浏览器将不会清除页面，而是阻止页面加载
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		// 第一次通过HTTPS请求，服务器响应Strict-Transport-Security
		// 浏览器后面尝试访问这个网站的请求都会自动把HTTP替换为HTTPS
		// 过期时间到了，后面通过HTTP的访问恢复到正常模式，不会再自动跳转到HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}

	// Also consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
}
