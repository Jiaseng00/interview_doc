package pkg

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"sync"
)

func DataHandler(c *gin.Context) {
	//// 调用低效计算
	dataSize := rand.Intn(1000) + 1000
	result := Calculation(dataSize)

	// 返回响应
	c.JSON(200, gin.H{
		"message": "Data processed",
		"result":  result,
	})
}

func Calculation(dataSize int) int {
	// 通过大量的无意义延时和计算模拟低效的计算
	var wg sync.WaitGroup
	resultChan := make(chan int, dataSize)

	// 并行计算每个数据
	for i := 0; i < dataSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resultChan <- ProcessData(i) // 并发执行计算
		}(i)
	}

	// 关闭 channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 汇总所有 goroutine 计算的结果
	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total
}

func ProcessData(x int) int {
	// 一个低效的计算方法，用来增加计算时间
	result := 0
	randValues := make([]int, 5000)
	for i := 0; i < 5000; i++ {
		randValues[i] = rand.Intn(100)
	}

	// 进行大量无用的循环，增加计算耗时
	for i := 0; i < 5000; i++ { // 扩大循环次数，增加 CPU 占用
		result += randValues[i]
	}

	return result + x
}
