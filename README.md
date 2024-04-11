# golang-reverse-proxy-middleware

#### 介绍

golang-reverse-proxy-middleware

#### 软件架构

软件架构说明

#### 安装教程

1. `cd static && pnpm install && pnpm run build && cd ..`

2. `cp static/dist/* public/`

3. `go build main.go`

#### 使用说明

```
go run main.go
```

设置环境变量 token 访问秘钥token123456

设置环境变量 port 监听端口8080

访问地址:

`http://localhost:8080/token/token123456/https/www.360.cn`

`http://localhost:8080/token/token123456/http/example.com`

# 设定代理行为的重定向方式

可以设定请求头中的字段"x-proxy-redirect"为"error" | "follow" |
"manual"来设定代理行为的重定向方式.
