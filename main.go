package main

import (
	// "net/http"

	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseForwardedHeader 解析 "Forwarded" HTTP 头部信息，返回一个 ForwardedBy 结构体切片。
// header: 代表被转发的请求的 "Forwarded" 头部字符串。
// 返回值: 一个包含所有转发标识的 ForwardedBy 结构体切片，以及可能发生的错误。
func parseForwardedHeader(header string) ([]ForwardedBy, error) {
	var forwardedByList []ForwardedBy
	parts := strings.Split(header, ", ")

	for _, part := range parts {
		for _, param := range strings.Split(part, ";") {
			param = strings.TrimSpace(param)
			if !strings.HasPrefix(param, "by=") {
				continue
			}

			// 分离 by 参数的值
			value := strings.TrimPrefix(param, "by=")
			// host, port, err := net.SplitHostPort(value)
			// if err != nil {
			// 如果没有端口信息，host 就是整个值
			var host = value
			// port = ""
			// }

			forwardedBy := ForwardedBy{
				Identifier: host,
				// Port:       port,
			}

			// 检查是否重复
			// isDuplicate := false
			// for _, existing := range forwardedByList {
			// 	if existing.Identifier == forwardedBy.Identifier && existing.Port == forwardedBy.Port {
			// 		isDuplicate = true
			// 		break
			// 	}
			// }
			// if !isDuplicate {
			forwardedByList = append(forwardedByList, forwardedBy)
			// }
		}
	}

	return forwardedByList, nil
}

func proxyHandler(w http.ResponseWriter, r *http.Request, url *url.URL) {
	fmt.Println("method:", r.Method)
	fmt.Println("url:", r.URL)
	fmt.Println("host:", r.Host)
	log.Println("proxyHandler", "header:")
	/*/* 这里删除除了第一次请求的 Proxy-Authorization  删除代理认证信息 */

	// r.Header.Del("Proxy-Authorization")
	clienthost, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("clienthost:", clienthost)
	log.Println("clientport:", port)
	forwarded := fmt.Sprintf(
		"for=%s;by=%s;host=%s;proto=%s",
		r.Header.Get("true-client-ip"), // clienthost, // 代理自己的标识或IP地址
		r.Host,                         // 代理的标识
		r.Host,                         // 原始请求的目标主机名
		"http",                         // 或者 "https" 根据实际协议
	)
	r.Header.Add("Forwarded", forwarded)
	for k, v := range r.Header {
		// fmt.Println("key:", k)
		log.Println("proxyHandler", k, ":", strings.Join(v, ","))
	}
	forwardedHeader := strings.Join(r.Header.Values("Forwarded"), ", ")
	log.Println("forwardedHeader:", forwardedHeader)
	forwardedByList, err := parseForwardedHeader(forwardedHeader)
	log.Println("forwardedByList:", forwardedByList)
	if len(forwardedByList) != len(setFromForwardedBy(forwardedByList)) {
		w.WriteHeader(508)
		fmt.Fprintln(w, "Duplicate 'by' identifiers found in 'Forwarded' header.")
		log.Println("Duplicate 'by' identifiers found in 'Forwarded' header.")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing 'Forwarded' header: %v", err)
		return
	}
	targetUrl := url.String()
	/*r.URL可能是http://开头,也可能只有路径  */
	// if startsWithHTTP(r.URL.String()) {
	// 	targetUrl = r.URL.String()
	// }
	// 这里假设目标服务器都是HTTP的，实际情况可能需要处理HTTPS
	fmt.Println("targetUrl:", targetUrl)
	// 创建一个使用了代理的客户端
	defer r.Body.Close()
	/* 请求body的问题 */
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println("body:", string(bodyBytes))
	client := &http.Client{ /* Transport: newTransport("http://your_proxy_address:port") */
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},

		/* Jar: jar */} // 替换为你的代理服务器地址和端口

	if r.Header.Get("x-proxy-redirect") == "follow" {
		client.CheckRedirect = nil
	}
	proxyReq, err := http.NewRequest(r.Method, targetUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header.Clone()

	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the response to the client's response.
	for k, v := range resp.Header {
		w.Header().Add(k, strings.Join(v, ","))
	}
	w.WriteHeader(resp.StatusCode)

	// Copy the response body back to the client.
	bodyBytes2, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(bodyBytes2)); err != nil {
		log.Println("Error writing response:", err)
	}
}

// 辅助函数：将ForwardedBy列表转换为集合（set），用于快速判断重复项
func setFromForwardedBy(forwardedByList []ForwardedBy) map[string]bool {
	set := make(map[string]bool)
	for _, fb := range forwardedByList {
		set[fb.Identifier] = true
	}
	return set
}

type ForwardedBy struct {
	Identifier string
}

func main() {
	var token string
	value, exists := os.LookupEnv("token")

	if exists {
		fmt.Println("环境变量 token 存在，其值为:", value)
		token = value
	} else {
		fmt.Println("环境变量 token 不存在或为空")
		token = "token123456" // 默认值
	}
	var port string
	value2, exists2 := os.LookupEnv("port")

	if exists2 {
		fmt.Println("环境变量 port 存在，其值为:", value2)
		port = value2
	} else {
		fmt.Println("环境变量 port 不存在或为空")
		port = "8080" // 默认值
	}
	// 创建一个 Gin 引擎实例
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {

		ctx.Writer.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		ctx.Next()
	}, func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/token/"+token+"/http/") {

			var url, err = url.Parse("http://" + ctx.Request.URL.Path[len("/token/"+token+"/http/"):])
			if err != nil {
				log.Println(err)
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			proxyHandler(ctx.Writer, ctx.Request, url)
			return
		} else if strings.HasPrefix(ctx.Request.URL.Path, "/token/"+token+"/https/") {
			var url, err = url.Parse("https://" + ctx.Request.URL.Path[len("/token/"+token+"/https/"):])
			if err != nil {
				log.Println(err)
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			proxyHandler(ctx.Writer, ctx.Request, url)
			return
		}
		ctx.Next()
	})
	// 设置静态文件服务
	r.Static("/", "./public")

	// 定义一个处理函数，用于渲染 "Hello, World!" 静态网页

	// 启动服务器并监听端口
	r.Run(":" + (port))
}
