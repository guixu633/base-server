package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/guixu633/base-server/module/config"
	"github.com/guixu633/base-server/pkg/http"
	"github.com/guixu633/base-server/pkg/service"
)

func main() {
	var cli struct {
		ConfigFile string ` type:"path" help:"配置文件的路径" default:"config.toml" short:"c"`
	}

	// 创建命令行解析器
	parser := kong.Parse(&cli)

	// 解析命令行参数
	_, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing command line: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(cli.ConfigFile)
	if err != nil {
		os.Exit(1)
	}

	svc := service.NewService(cfg)
	if err != nil {
		os.Exit(1)
	}
	httpServer := http.NewServer(8001, svc, "certificate/guixuu.com.crt", "certificate/guixuu.com.key")
	if err != nil {
		os.Exit(1)
	}
	httpServer.ListenAndServe()
}
