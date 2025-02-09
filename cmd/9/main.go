package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Client struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	jsonData := `{
		"id": 123,
		"name": "张三",
		"email": "zhangsan@example.com",
		"created_at": "2024-02-04T15:04:05Z"
	}`
	var client Client
	err := json.Unmarshal([]byte(jsonData), &client)
	if err != nil {
		fmt.Println("json unmarshal err:", err)
		return
	}
	if client.CreatedAt.IsZero() {
		client.CreatedAt = time.Now()
	}
	result, err := json.Marshal(client)
	if err != nil {
		fmt.Println("json marshal err:", err)
		return
	}
	fmt.Println(string(result))
}
