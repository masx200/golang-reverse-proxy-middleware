# golang-reverse-proxy-middleware

#### 介绍
golang-reverse-proxy-middleware

#### 软件架构
软件架构说明


#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明


设置环境变量 token 访问秘钥token123456

设置环境变量 port 监听端口8080

访问地址:

`http://localhost:3000/token/token123456/https/www.360.cn`

`http://localhost:3000/token/token123456/http/example.com`

# 设定代理行为的重定向方式

可以设定请求头中的字段"x-proxy-redirect"为"error" | "follow" | "manual"来设定代理行为的重定向方式.