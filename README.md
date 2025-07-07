# simplejrpc-go

`simplejrpc-go` 是一个基于 Go 语言的 JSON-RPC 框架，提供了配置管理、日志记录、数据验证、国际化支持等功能，旨在帮助开发者快速搭建高性能、可维护的 JSON-RPC 服务器。

## 特性
- **配置管理**：支持从 JSON 文件加载配置数据，支持嵌套配置项的访问，可通过点符号访问。
- **日志记录**：基于 `go.uber.org/zap` 实现日志记录，支持日志文件旋转、输出到标准输出、设置日志级别等功能。
- **数据验证**：提供结构体字段验证功能，支持多种验证规则，如必填字段、最小长度、数值范围等。
- **国际化支持**：支持从文件加载国际化资源，目前仅支持 INI 文件格式。
- **网络通信**：实现了 JSON-RPC 服务器和请求处理功能，支持中间件机制。

## 安装
```sh
go get simplejrpc-go
```

## 配置管理

配置文件 config.json 示例：
```json
{
    "test": {
        "version": "1.0.0",
        "jsonrpc": {
            "sockets": "rpc.sock"
        },
        "logger": {
            "path": "logs/",
            "file": "{Y-m-d}.log",
            "level": "error",
            "stdout": false,
            "StStatus": 0,
            "rotateBackupLimit": 7,
            "writerColorEnable": true,
            "RotateBackupCompress": 9,
            "rotateExpire": "1d",
            "Flag": 44
        }
    },
    "prod": {
        "version": "1.0.0",
        "jsonrpc": {
            "sockets": "rpc.sock"
        },
        "logger": {
            "path": "logs/",
            "file": "{Y-m-d}.log",
            "level": "error",
            "stdout": false,
            "StStatus": 0,
            "rotateBackupLimit": 7,
            "writerColorEnable": true,
            "RotateBackupCompress": 9,
            "rotateExpire": "1d",
            "Flag": 44
        }
    }
}
```


项目内部读取配置文件示例：
```go
package main


import (
	"fmt"
	"path/filepath"

	"github.com/DemonZack/simplejrpc-go/core"
	"github.com/DemonZack/simplejrpc-go/core/config"
	"github.com/DemonZack/simplejrpc-go/os/gpath"
)

func main() {
	env := "test"
	fullPath, _ := filepath.Abs(filepath.Dir("."))
	gpath.GmCfgPath = filepath.Join(filepath.Dir(fullPath), "..")
	core.InitContainer(config.WithConfigEnvFormatterOptionFunc(env))

	val, err := core.Container.CfgFmt().GetValue("logger.level").String()
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] logger.level : ", val)

	val, err = core.Container.CfgFmt().GetValue("jsonrpc.sockets").String()
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] jsonrpc.sockets : ", val)
}

```

## 日志记录
使用 glog 模块进行日志记录，支持日志文件旋转、输出到标准输出、设置日志级别等功能。以下是日志记录示例
```go
package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/DemonZack/simplejrpc-go/core/glog"
)

func main() {
	m := map[string]any{
		"path":                 "logs/",
		"file":                 "{Y-m-d}.log",
		"level":                "error",
		"stdout":               false,
		"StStatus":             0,
		"rotateBackupLimit":    7,
		"writerColorEnable":    true,
		"RotateBackupCompress": 9,
		"rotateExpire":         "1d",
		"Flag":                 44,
	}

	// config, err := LoadConfig("./testdata/config.json")
	config, err := glog.LoadConfig(m)
	if err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	// 初始化日志
	logger, err := glog.NewLogger(config)
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	defer logger.Sync() // 刷新缓冲区的日志

	// 使用示例
	logger.Info("Logger initialized successfully",
		zap.String("path", config.Path),
		zap.String("file", config.File),
		zap.Int("backupLimit", config.RotateBackupLimit),
	)

	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// 结构化日志
	logger.Info("User logged in",
		zap.String("username", "john"),
		zap.Int("attempt", 3),
		zap.Duration("duration", time.Second*5),
	)

	zap.L().Info("ddddddddddddddddddddddddd")
}

```

项目内部日志使用示例：
配置好日志配置
```go
msg := "log record content"
core.Container.Log().Info(msg)
core.Container.Log().Debug(msg)
core.Container.Log().Error(msg)
```

## 数据验证
使用 gvalid 模块进行数据验证，支持多种验证规则，可自定义验证器。以下是数据验证示例：
```go
package main

import (
	"fmt"
	"path/filepath"

	"github.com/DemonZack/simplejrpc-go/core"
	"github.com/DemonZack/simplejrpc-go/core/config"
	"github.com/DemonZack/simplejrpc-go/os/gpath"
)

type ExampleUser struct {
	Username string `validate:"min_length:6#The length is too small"`
	Age      any    `validate:"required#Required parameters are missing Age|range:18,100|int#Test verification error return"`
	Email    string `validate:"required#Email address is required"`
}

func main() {
	env := "test"
	fullPath, _ := filepath.Abs(filepath.Dir("."))
	gpath.GmCfgPath = filepath.Join(filepath.Dir(fullPath), "..")
	core.InitContainer(config.WithConfigEnvFormatterOptionFunc(env))

	user := ExampleUser{
		Username: "test",
		Age:      15,
		Email:    "",
	}

	err := core.Container.Valid().Walk(&user)
	if err != nil {
		fmt.Println("Validation error:", err)
	}
}

```

## 国际化支持
使用 gi18n 模块进行国际化支持，支持从文件加载国际化资源，目前仅支持 INI 文件格式。以下是国际化支持示例：
```go
package main

import (
    "path/filepath"
    "github.com/DemonZack/simplejrpc-go/core/gi18n"
)

func main() {
    fPath, _ := filepath.Abs(filepath.Dir(""))
    path := filepath.Join(fPath, "testdata")
    gi18n.Instance().SetPath(path)
    gi18n.Instance().SetLanguage(gi18n.English.String())

    key := "Welcome"
    val := gi18n.Instance().T(key)
    println(val)

    gi18n.Instance().SetLanguage(gi18n.SimplifiedChinese.String())
    val = gi18n.Instance().T(key)
    println(val)
}
```

## 启动 JSON-RPC 服务器
使用 server 模块启动 JSON-RPC 服务器，支持中间件机制。以下是启动服务器示例：
```go
package main

import (
	"github.com/DemonZack/simplejrpc-go/net/gsock"
	"github.com/DemonZack/simplejrpc-go/server"
)

type CustomHandler struct{}

func (c *CustomHandler) Hello(req *gsock.Request) (any, error) {
	return "Hello World", nil
}

func (c *CustomHandler) ProcessRequest(req *gsock.Request) {
	println("[*] ProcessRequest. ", req)
}

func (c *CustomHandler) ProcessResponse(resp any) (any, error) {
	println("[*] ProcessResponse. ", resp)
	return resp, nil
}

func main() {
	mockSockPath := "zack.sock"

	ds := server.NewDefaultServer(
		gsock.WithJsonRpcSimpleServiceHandler(gsock.NewJsonRpcSimpleServiceHandler()),
		gsock.WithJsonRpcSimpleServiceMiddlewares([]gsock.RPCMiddleware{
			&CustomMiddleware{},
		}),
	)

	hand := &CustomHandler{}
	ds.RegisterHandle("hello", hand.Hello, []gsock.RPCMiddleware{hand}...)
	err := ds.StartServer(mockSockPath)
	if err != nil {
		panic(err)
	}
}

type CustomMiddleware struct{}

func (c *CustomMiddleware) ProcessRequest(req *gsock.Request) {
	println("[*] ProcessRequest. ", req)
}

func (c *CustomMiddleware) ProcessResponse(resp any) (any, error) {
	println("[*] ProcessResponse. ", resp)
	return resp, nil
}

```

## 贡献
欢迎贡献代码，请遵循以下步骤：

1. Fork 仓库
2. 创建新的分支
3. 提交代码并推送至新分支
4. 发起 Pull Request