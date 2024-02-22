# FileWeb
轻量级文件映射小工具,用于临时映射本地目录到Web服务器,可在局域网之间下载文件使用。
- 自定义映射本地路径为虚拟路径
- 虚拟路径支持密码访问

> 首次启动后会自动生成配置文件

```yaml
#服务启动端口
port: 10261
#是否开启目录索引
folderIndex: true
#用户列表(修改后直接生效,无需重启)
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
```