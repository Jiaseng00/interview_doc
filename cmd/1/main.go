package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	taskQueue := make(chan int, 3)
	var wg sync.WaitGroup
	numWorkers := 3 //模拟3个员工

	// 每有一个员工，加一个Goroutine
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskQueue, &wg)
	}

	// 给予10个任务给这些员工
	for i := 1; i <= 10; i++ {
		taskQueue <- i
		fmt.Printf("任务%d加入队列\n", i)
	}
	close(taskQueue)

	wg.Wait()
	fmt.Println("所有任务完成")
}

func worker(id int, taskQueue chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskQueue {
		fmt.Printf("Worker %d 正在处理 %d\n", id, task)
		time.Sleep(time.Second)
	}
}
