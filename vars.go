package main

var (
	FileWebVersion   = "0.0.1"
	DefaultConfigTpl = `
port: 10261 #服务启动端口
#用户列表
user:
  - name: "root"
    password: "123"

#映射路径,修改映射目录需要重新启动服务
mapping:
  - path: "/" #虚拟路径,例如访问 http://127.0.0.1/ 会直接访问到配置的Local目录下
    local: "./" #本地映射路径
    auth: false #目录是否需要认证,如果为true则会需要user的认证信息
  - path: "/auth" #,例如访问 http://127.0.0.1/auth 会直接访问到配置的Local目录下
    local: "./"
    auth: true
`
)
