package main

type User struct {
	Name     string
	Password string
}
type Mapping struct {
	Path  string //虚拟映射路径
	Local string //本地映射路径
	Auth  bool   //目录是否需要认证
}
type Config struct {
	Port    uint16 //端口
	User    []User
	Mapping []Mapping
}
