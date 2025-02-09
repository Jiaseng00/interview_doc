package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkDataHandler(b *testing.B) {
	r := gin.Default()

	// 将 / 路由绑定到 DataHandler 函数
	r.GET("/", DataHandler)

	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Error creating request: %v", err)
	}

	// 创建一个模拟的 HTTP 响应记录器
	rr := httptest.NewRecorder()

	// 重置计时器，以避免初始化开销
	b.ResetTimer()

	// 执行基准测试，模拟请求处理的多次执行
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(rr, req) // 调用 DataHandler
	}
}

func BenchmarkCalculation(b *testing.B) {
	dataSize := 1000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Calculation(dataSize)
	}
}

func BenchmarkProcessData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ProcessData(i)
	}
}
