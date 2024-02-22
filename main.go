package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/table"
	"strings"

	"github.com/spf13/viper"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var config *Config

func BasicAuth(handler http.Handler, auth bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if auth {
			isAuth := false
			for _, u := range config.User {
				if u.Name == user && pwd == u.Password && ok {
					isAuth = true
					break
				}
			}
			if !isAuth {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{}
	loadConfig("config.yaml", func(ctx context.Context) {
		printUser()
	})
	addr := fmt.Sprintf(":%d", config.Port)
	server.Addr = addr

	for _, mp := range config.Mapping {
		fs := FileSystem{
			fs: http.Dir(mp.Local),
		}
		filesServer := http.FileServer(fs)
		path := mp.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		mux.Handle(path, BasicAuth(http.StripPrefix(mp.Path, filesServer), mp.Auth))
	}
	printMapping()
	fmt.Printf("FileWeb Version: %s\n", FileWebVersion)
	fmt.Printf("FileWeb 服务启动在 http://127.0.0.1%s\n", addr)
	server.Handler = mux
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(fmt.Sprintf("FileWeb 启动失败: %s", err.Error()))
	}
}

func loadConfig(name string, reloadEvent func(ctx context.Context)) {
	ctx := context.Background()
	viper := viper.New()
	viper.SetConfigFile(name)
	if err := viper.ReadInConfig(); err != nil {
		if os.IsExist(err) {
			log.Fatalln(fmt.Sprintf("[配置文件] [%s]加载失败: %s", name, err.Error()))
		} else {
			if err := os.WriteFile(name, []byte(DefaultConfigTpl), 0777); err != nil {
				log.Fatalln(fmt.Sprintf("[配置文件] 初始化配置文件失败: %s", err.Error()))
			} else {
				log.Println(fmt.Sprintf("[配置文件] 生成默认配置文件成功[%s],请修改后尝试重启服务", name))
				os.Exit(0)
			}

		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&config); err != nil {
			log.Println(fmt.Sprintf("[%s]配置文件重载失败: %s", in.Name, err.Error()))
		} else {
			log.Println(fmt.Sprintf("[%s]配置文件重载完成", in.Name))
			reloadEvent(ctx)
		}
	})
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(fmt.Sprintf("[%s]解析配置文件失败: %s", name, err.Error()))
	}
	log.Println(fmt.Sprintf("[%s]配置文件加载完成", name))
	reloadEvent(ctx)
}

func printUser() {
	t, err := gotable.Create("用户名", "密码")
	if err != nil {
		log.Fatalln(fmt.Sprintf("初始化用户表格失败: %s", err.Error()))
	}
	values := make([]map[string]string, 0)
	for _, u := range config.User {
		values = append(values, map[string]string{
			"用户名": u.Name,
			"密码":  u.Password,
		})
	}
	t.Align("用户名", table.L)
	t.Align("密码", table.L)
	t.AddRows(values)
	fmt.Println(t)
}

func printMapping() {
	t, err := gotable.Create("虚拟目录", "本地路径", "认证状态")
	if err != nil {
		log.Fatalln(fmt.Sprintf("初始化映射表格失败: %s", err.Error()))
	}
	values := make([]map[string]string, 0)
	for _, u := range config.Mapping {
		authText := "False"
		if u.Auth {
			authText = "True"
		}
		values = append(values, map[string]string{
			"虚拟目录": u.Path,
			"本地路径": u.Local,
			"认证状态": authText,
		})
	}
	t.Align("虚拟目录", table.L)
	t.Align("本地路径", table.L)
	t.Align("认证状态", table.L)
	t.AddRows(values)
	fmt.Println(t)
}

type FileSystem struct {
	fs               http.FileSystem
	readDirBatchSize int
}

func (fs FileSystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return ReadDirFile{f, 2}, nil
}

type ReadDirFile struct {
	http.File
	readDirBatchSize int
}

func (f ReadDirFile) Stat() (fs.FileInfo, error) {
	s, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
	LOOP:
		for {
			fl, err := f.File.Readdir(f.readDirBatchSize)
			switch err {
			case io.EOF:
				break LOOP
			case nil:
				for _, f := range fl {
					if f.Name() == "index.html" {
						return s, err
					}
				}
			default:
				return nil, err
			}
		}
		return nil, os.ErrNotExist
	}
	return s, err
}
