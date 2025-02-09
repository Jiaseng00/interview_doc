package main

import (
	"Good_Net/cmd/8/pkg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof" // 启动 pprof 监控
)

func main() {
	go func() {
		// 启动 pprof 服务，监听 6060 端口
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// 创建一个 Gin 路由实例
	r := gin.Default()

	// 定义路由，调用 handler 包中的 ProcessData 函数
	r.GET("/", pkg.DataHandler)

	// 启动主 HTTP 服务器，监听 8080 端口
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
