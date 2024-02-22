package main

type Config struct {
	Port uint16 //端口
	User []struct {
		Name     string
		Password string
	}
	Mapping []struct {
		Path  string //虚拟映射路径
		Local string //本地映射路径
		Auth  bool   //目录是否需要认证
	}
}
