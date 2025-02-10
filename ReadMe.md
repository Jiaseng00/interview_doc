### 1. Golang并发模型
***
**Goroutine**
+ 作用：是个能在后台轻量运行并且不会干扰到主程序运行的线程。它允许程序同时执行多个逻辑，而且不会阻挡主程序
+ 原理：由Go runtime管理，轻型并动态的线程。Goroutine会被调度到多个OS线程上执行，并且自动切换他们的执行

**Channel**
+ 作用：是个重要工具让Goroutine来传送之间的数据，实现不同goroutine之间的数据通信
+ 原理：Channel是类型化的，可以传递特定类型的数据。具有阻塞特性，发送和接收数据会阻塞直到对方准备好。
缓冲Channel可以避免阻塞.
### Go 代码
```go
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
```
显示
```bash
任务1加入队列
任务2加入队列
任务3加入队列
任务4加入队列
任务5加入队列
任务6加入队列
Worker 0 正在处理 1
Worker 2 正在处理 2
Worker 1 正在处理 3
Worker 1 正在处理 4
Worker 0 正在处理 6
任务7加入队列
任务8加入队列
任务9加入队列
Worker 2 正在处理 5
Worker 0 正在处理 7
Worker 1 正在处理 8
任务10加入队列
Worker 2 正在处理 9
Worker 1 正在处理 10
所有任务完成
```
### 2. 数据库设计与查询优化
___
查询最近7天注册的用户<br>
[GORM](https://gorm.io/) 安装GORM包链接MySQL
```
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```
Go 代码
```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	Id        int    `gorm:"id" json:"id"`
	Name      string `gorm:"name" json:"name"`
	Email     string `gorm:"email" json:"email"`
	CreatedAt string `gorm:"created_at" json:"created_at"`
}

func main() {
	// 模拟链接MySQL数据库
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}
	// 链接完毕
	var users []User
	sevenDaysAgoTime := time.Now().AddDate(0, 0, -7).Unix()
	err = db.Table("users").Where("created_at >= ?", sevenDaysAgoTime).Find(&users).Error
	if err != nil {
		log.Println(err)
		return
	}
}
```
2. 为了优化查询性能，必须在MySQL的数据库里，用户表(users)里设定Index来优化搜索的速度。
例如email必须是不允许重复的，可以让email成为一个unique Index，来让数据库搜索减少不必要的步骤。

### 3. RESTful API 设计
***
##### 订单管理API
服务器运行后,显示订单CRUD Api
```bash
[GIN-debug] GET    /                         --> main.GetOrder (3 handlers)
[GIN-debug] POST   /                         --> main.CreateOrder (3 handlers)
[GIN-debug] PUT    /:id                      --> main.UpdateOrder (3 handlers)
[GIN-debug] DELETE /:id                      --> main.DeleteOrder (3 handlers)
```
订单API路由结构
```go
package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()

	r.Group("/api/orders")
	{
		r.GET("", GetOrder)
		r.POST("/", CreateOrder)
		r.PUT("/:id", UpdateOrder)
		r.DELETE("/:id", DeleteOrder)
	}

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
```
JSON响应助手，简化回复前端的逻辑代码
```go
// JSON 响应助手
func JSONResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
		"data":    nil,
	})
}
```
创建订单的API代码
```go
// 模拟数据库里Order的构造
type Order struct {
	Id         int `gorm:"id" json:"id"`
	CustomerId int `gorm:"customer_id" json:"customer_id"`
	Status     int `gorm:"status" json:"status"`
	CreatedAt  int `gorm:"created_at" json:"created_at"`
	UpdatedAt  int `gorm:"updated_at" json:"updated_at"`
}

func CreateOrder(c *gin.Context) {
	var db *gorm.DB //模拟已经链接数据库的情况
	var order Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "无效的JSON格式")
		return
	}
	err = db.Table("orders").Create(&order).Error
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建订单失败")
		return
	}
	JSONResponse(c, http.StatusCreated, "成功创建订单", order)
}
```
### 4. JWT⾝份验证
***
JWT的基本构造是由`header.payload.signature`链接和加密组成而生成的一个token.
+ **签名**是把header和payload经过Base64Url编码后，加上一个开发员独立收藏的密钥包在一起，
然后使用其中一个算法来进行签名，最后连接在JWT`signature`的部分。
```bash
// 签名
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```
+ **解析**会将Base64Url编码的Header和Payload部分解码，提取其中的JSON数据。查看JWT中的
用户信息或其他声明的步骤。
+ **验证**是当服务器接收到一个请求发来的token,服务器必须使用相同的签名和密钥对header和
payload进行重新签名。如果请求发来的签名和服务器重新生成的签名一致，说明JWT没有被篡改，可以
进行下个程序行动。为了JWT的有效性，开发员通常会设置有效期，当JWT无效，用户必须再次
登录来获得新的JWT。

```go
package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func main() {
	// 用户登入生成JWT
	token, err := GenerateToken("example123")
	if err != nil {
		log.Fatal(err)
	}
	ParseToken(token)
}

// JWT生成
func GenerateToken(username string) (string, error) {
	// 模拟用户请求登入Api，回车需要的资料
	// JWT Payload资料生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 放入需要的用户资料
		"foo":      "bar",
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	// 从.env文件，调用密钥来进行签名
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tokenString)
	return tokenString, nil
}

// JWT验证
func ParseToken(tokenString string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := []byte(os.Getenv("JWT_SECRET"))

		return secret, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}

```
### 5. Golang项⽬结构与代码规范
***
```
project_folder
├───.idea               #GoLand的IDE配置文件夹
├───config              #存放应用配置文件
├───app                 #存放应用的核心代码
│   ├───controllers     #控制器：处理请求和响应的逻辑
│   ├───models          #模型：定义数据结构
│   ├───repositories    #数据仓库：与数据库进行交互沟通
│   └───services        #服务：业务逻辑处理，协调不同层的操作
├───routes              #路由：定义HTTP路由和请求处理的映射
├───middleware          #中间件：处理跨切面功能，如JWT认证
└───utils               #工具：存放辅助函数和常用工具代码
```
+ **config** 存档项目配置文件(如`YAML`或`JSON`)。包括数据库链接，应用接口，
API秘钥等。
+ **app** 存放应用核心逻辑。
  + **controller** 负责HTTP请求和响应处理
  + **models** 存放数据结构和数据库模型定义
  + **repositories** 负责与数据库进行交接的逻辑操作
  + **services** 负责业务逻辑，协调控制器和仓库之间的操作
+ **routes** 定义HTTP路由。将路由和控制器的处理绑定
+ **middleware** 存放中间件，如JWT
+ **utils** 存放常用的工具函数。例如JSON响应助手
___
配置管理
+ 为了方便管理一些固定，敏感，斌且重要的配置像是数据库链接配置，API_Key，缓存链接
配置等等，可以把这些重要资料卸载`.env`文件里(使用gitignore忽略此文件)。然后使用
像是`https://github.com/joho/godotenv.git` 的包来从`.env`文件调取需要的配置

**.env文件**
```
API_KEY=APIKEYEXAMPLE
PORT=1234
```
**读取.env文件中的函数**
```go
package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  apiKey := os.Getenv("API_KEY")
  port := os.Getenv("PORT")

  //继续主程序逻辑
}
```
日志管理
+ 为了良好的扩充性，使用日志轮转和保留可以避免信息的遗失并且能不断储存新的信息。
  + 此例子使用`logrus`和`lumberjack`来进行日志管理

日志文件轮转和保留的设定
```go
fileLogger := &lumberjack.Logger{
    Filename:   "service.log", // 日志文件的路径和名称
    MaxSize:    100,           // 文件最大大小，单位 MB，超过该大小会进行轮转
    MaxBackups: 30,            // 保留的旧日志文件最大数量，超过会被删除
    MaxAge:     28,            // 保留日志的天数，超过会被删除
    Compress:   true,          // 是否压缩旧的日志文件
}
```
在进行信息保存前，可以把信息以JSON格式存入日志文件，以便需要使用第三方软件来导入或导出信息。

service.log里的记录
```
{"level":"info","msg":"Service started","time":"2025-02-07T23:23:10+08:00"}
{"level":"warning","msg":"TEST","time":"2025-02-07T23:23:10+08:00"}
{"level":"info","msg":"Handling request","request_id":"12345","service_name":"UserService","time":"2025-02-07T23:23:10+08:00"}
{"error":"Database connection failed","level":"error","msg":"Failed to process request","time":"2025-02-07T23:23:10+08:00"}
```
在开发阶段，推荐把日志输出在控制台并使用颜色代码，以便阅读
```bash
INFO[70717-07+08 23:711:27] Service started                              
WARN[70717-07+08 23:711:27] TEST
INFO[70717-07+08 23:711:27] Handling request                              request_id=12345 service_name=UserService
ERRO[70717-07+08 23:711:27] Failed to process request                     error="Database connection failed"
```
错误处理
+ 处理简单的错误，可以直接使用`fmt.Println(error)`或者`log.Println(error)`
显示错误。但为了良好的扩充性，可以使用**错误包装**技巧来处理错误，这会让开发员更加容
易的了解问题的所在.例如：
```go
var ErrNotFound = errors.New("resource not found")

func someFunction() error {
    return fmt.Errorf("failed to get resource: %w", ErrNotFound)
}

func main() {
    err := someFunction()
    if errors.Is(err, ErrNotFound) {
        fmt.Println("Resource not found!")
    }
}
```
数据库连接
+ 可以使用ORM框架(如`GORM`)和调用`.env`文件里的数据库配置来简化链接数据库的操作，
并且使用连接池优化与数据库的链接
```go
import (
  "database/sql"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "os"
)

  var (
      dbUser     = os.Getenv("DB_USER")
      dbPassword = os.Getenv("DB_PASSWORD")
      dbHost     = os.Getenv("DB_HOST")
      dbPort     = os.Getenv("DB_PORT")
      dbName     = os.Getenv("DB_NAME")
  )

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s",
        dbUser, dbPassword, dbHost, dbPort, dbName, mysqlTLSName)

sqlDB, err := sql.Open("mysql", dsn)

sqlDB.SetMaxOpenConns(10)
sqlDB.SetMaxIdleConns(5)

gormDB, err := gorm.Open(mysql.New(mysql.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```

### 6. Golang错误处理与⽇志管理
***
1. `errors.New`是能让开发员写出单个简单的错误类型的信息。
2. `fmt.Errorf`是能够使用`%w`来显示出所收到的`err`错误类代码，并且能够自行加入需要的信息，和`fmt.Printf`类似。
3. `errors.Wrap`是让开发员能够把一个错误配合一个信息结合一起发送的工具，并且显示堆栈信息，有需要时能够用`errors.Unwrap`取出错误信息

Go代码
```go
package main

import (
  "errors"
  "fmt"
  "github.com/natefinch/lumberjack"
  pkgErr "github.com/pkg/errors"
  "github.com/sirupsen/logrus"
  "os"
  "time"
)

var (
  ErrorDB    = errors.New("DB connection error")
  ErrorTwo   = errors.New("error two")
  ErrorThree = errors.New("error three")
)

func main() {
  log := logrus.New()
  // Set up file-based log rotation with JSON format
  fileLogger := &lumberjack.Logger{
    Filename:   "service.log", // Log file name
    MaxSize:    100,           // Max size in MB before rotating
    MaxBackups: 30,            // Max number of old log files to keep
    MaxAge:     28,            // Max number of days to retain old log files
    Compress:   true,          // Compress old log files
  }

  // Create a separate logger for file logging with JSON format
  fileLog := logrus.New()
  fileLog.SetFormatter(&logrus.JSONFormatter{})
  fileLog.SetOutput(fileLogger)

  // Create a separate logger for console logging with color
  consoleLog := logrus.New()
  consoleLog.SetFormatter(&logrus.TextFormatter{
    ForceColors:     true, // Forces colorization on the output
    FullTimestamp:   true, // Adds full timestamps
    TimestampFormat: fmt.Sprintf("%s", time.Now().UTC().Format("2006-01-02 15:04:05")),
  })
  consoleLog.SetOutput(os.Stdout)

  // Set log level (Info and above)
  log.SetLevel(logrus.InfoLevel)

  // Log to console and file
  consoleLog.Info("Service started")
  fileLog.Info("Service started")

  consoleLog.Warn("TEST")
  fileLog.Warn("TEST")

  // Log with contextual information
  consoleLog.WithFields(logrus.Fields{
    "service_name": "UserService",
    "request_id":   "12345",
  }).Info("Handling request")

  fileLog.WithFields(logrus.Fields{
    "service_name": "UserService",
    "request_id":   "12345",
  }).Info("Handling request")

  // Simulate an error
  err := someFunction(ErrorTwo)
  if err != nil {
    // Use pkgErr.WithStack to add stack trace information
    wrappedErr := pkgErr.Wrap(ErrorDB, "something wrong in DB")
    // Log the error with stack information
    consoleLog.WithFields(logrus.Fields{
      "error": wrappedErr,
    }).Error("Failed to process request")

    fileLog.WithFields(logrus.Fields{
      "error": wrappedErr,
    }).Error("Failed to process request")
  }

}

func someFunction(err error) error {
  return fmt.Errorf("error : %w", err)
}

```
终端显示的logger
```bash
INFO[90933-09-09 08:02:44] Service started
WARN[90933-09-09 08:02:44] TEST
INFO[90933-09-09 08:02:44] Handling request                              request_id=12345 service_name=UserService
ERRO[90933-09-09 08:02:44] Failed to process request                     error="something wrong in DB: DB connection error"
```
以JSON格式记录错误
```
{"level":"info","msg":"Service started","time":"2025-02-09T16:10:26+08:00"}
{"level":"warning","msg":"TEST","time":"2025-02-09T16:10:26+08:00"}
{"level":"info","msg":"Handling request","request_id":"12345","service_name":"UserService","time":"2025-02-09T16:10:26+08:00"}
{"error":"something wrong in DB: DB connection error","level":"error","msg":"Failed to process request","time":"2025-02-09T16:10:26+08:00"}
```

### 7. API并发安全问题
***
在高并发场景下，可能会导致前一秒才被更新的数据被后一秒的更行覆盖，导致数据的损失。<br>
以钱包交易作为例子：
+ 当用户**A**和用户**B**都请求对同个数据进行更新，两边刚好**同时**发送请求，形成**并发式**请求。
+ 数据A的余额时**100元**，用户A和用户B**同时收到**数据A的余额是100元
+ 用户A请求增加余额(+50元)，数据的余额更新为150元
+ 用户B请求增加余额(+100元)，由于用户B的目前获得的数据余额还是100元，数据的余额更新为200元，覆盖用户A更新后的数据
+ 其实真正的余额应该是100+50+100,是**250元**
+ 数据不一致，造成**数据流失**

从这例子可以看出，当在高并发场景下，数据的一致性的重要性。

解决方法
1. 使用事务
+ 事务(Transaction)确保所有的线程成功进行，才能允许数据库进行修改，如有问题，放弃所有线程并且回车到开始线程之前并不会对数据库造成改变。
+ 这方法只能一定程度上保证数据的一致性，比较类似于当有什么线程失败时，让所有数据改动回到原点，以免造成数据的不一致性。

2. 乐观锁
+ 在数据库的表里加入`version`来监督数据的更新动向
+ 当进行更新线程时，比较程序和数据库的`version`，如果不一致就取消更新，一致的话就进行更新。
+ 在高并发场景和事务一起使用，保证良好的数据一致性，避免数据丢失或冲突

3. 悲观锁
+ 在进行任何动作前，使用`for update`锁来确保数据交接进行期间，只有此线程正在进行，并且其他的线程必须等待此锁打开后，再次获得此锁才能进行。

```go
package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	var db *gorm.DB //模拟数据库链接
	err := updateBalanceWithPessimisticLock(db, 1, 50, "deposit")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = updateBalanceWithOptimisticLock(db, 1, 100, "deposit")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 悲观锁
func updateBalanceWithPessimisticLock(db *gorm.DB, userID int, amount float64, action string) error {
	type wallet struct {
		ID      int     `gorm:"id"`
		UserId  int     `gorm:"user_id"`
		Balance float64 `gorm:"balance"`
	}
	var userWallet wallet
	tx := db.Begin()

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).Find(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	switch action {
	case "deposit":
		userWallet.Balance += amount
	case "withdrawal":
		if (userWallet.Balance - amount) >= 0 {
			userWallet.Balance -= amount
		} else {
			tx.Rollback()
			return errors.New("insufficient balance ")
		}
	default:
		tx.Rollback()
		return errors.New("invalid action,only \"deposit\" and \"withdrawal\" are supported")
	}

	err = tx.Model(&wallet{}).Where("user_id = ?", userID).Updates(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 乐观锁
func updateBalanceWithOptimisticLock(db *gorm.DB, userID int, amount float64, action string) error {
	type wallet struct {
		ID      int     `gorm:"id"`
		UserId  int     `gorm:"user_id"`
		Balance float64 `gorm:"balance"`
		Version int     `gorm:"version"`
	}
	var userWallet wallet
	tx := db.Begin()

	err := tx.Model(&wallet{}).Where("user_id = ?", userID).Find(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	switch action {
	case "deposit":
		userWallet.Balance += amount
	case "withdrawal":
		if (userWallet.Balance - amount) >= 0 {
			userWallet.Balance -= amount
		} else {
			tx.Rollback()
			return errors.New("insufficient balance ")
		}
	default:
		tx.Rollback()
		return errors.New("invalid action,only \"deposit\" and \"withdrawal\" are supported")
	}

	err = tx.Model(&wallet{}).Where("user_id = ? AND version = ?", userID, userWallet.Version).
		Updates(&wallet{Balance: userWallet.Balance, Version: userWallet.Version + 1}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
```
### 8.代码性能优化
***
为了排查性能问题：
1. 可以在服务器上开启另一个不同的端口，专门拿来做性能测试，这样做就不需要担心会太大的影响到主程序端口的性能
。这个API回应时间是15秒左右
```go
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

```
2. 在终端输入`go tool pprof http://localhost:6060/debug/pprof/profile
   `来收集30秒内（默认）CPU的profile。然后选用如`web`,`top`,`list` 等来查看占内存高的函数。
3. 收集CPU Profile后，输入`top`来查看占用CPU性能比较高的函数有哪些，从这里可以看见`pkg.ProcessData`,`rand.(*Rand).Intn`使用的CPU内存
是非常高的。
```
Showing nodes accounting for 1060ms, 97.25% of 1090ms total
Showing top 10 nodes out of 56
      flat  flat%   sum%        cum   cum%
     420ms 38.53% 38.53%      720ms 66.06%  math/rand.(*Rand).Int31n
     150ms 13.76% 52.29%      150ms 13.76%  internal/chacha8rand.block
     140ms 12.84% 65.14%      860ms 78.90%  math/rand.(*Rand).Intn
     110ms 10.09% 75.23%     1030ms 94.50%  math/rand.Intn
      90ms  8.26% 83.49%      250ms 22.94%  runtime.rand
      40ms  3.67% 87.16%      290ms 26.61%  math/rand.(*runtimeSource).Int63
      40ms  3.67% 90.83%       60ms  5.50%  math/rand.globalRand
      30ms  2.75% 93.58%       30ms  2.75%  runtime.cgocall
      20ms  1.83% 95.41%     1050ms 96.33%  Good_Net/cmd/8/pkg.ProcessData
      20ms  1.83% 97.25%       20ms  1.83%  sync/atomic.(*Pointer[go.shape.struct { math/rand.src math/rand.Source; math/rand.s64 math/rand.Source64; math/rand.readVal int64; math/rand.readPos int8 }]).Load (inline)
```
4.  调用`list`指令来查看ProcessData的逻辑代码。
```
Total: 1.09s
ROUTINE ======================== Good_Net/cmd/8/pkg.ProcessData in D:\Good_Net\cmd\8\pkg\pkg.go
      20ms      1.05s (flat, cum) 96.33% of Total
         .          .     37:func ProcessData(x int) int {
         .          .     38:   // 一个低效的计算方法，用来增加计算时间
         .          .     39:   result := 0
         .          .     40:   // 进行大量无用的循环，增加计算耗时
      20ms       20ms     41:   for i := 0; i < 5000; i++ { // 扩大循环次数，增加 CPU 占用
         .      1.03s     42:           result += rand.Intn(100)
         .          .     43:   }
         .          .     44:   return result + x
         .          .     45:}
```
5.  使用`Benchmark`工具来写出函数测试文件`pkg_test.go`。
```go
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

func BenchmarkInefficientCalculation(b *testing.B) {
	dataSize := 1000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InefficientCalculation(dataSize)
	}
}

func BenchmarkProcessData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ProcessData(i)
	}
}
```
##### 优化前
6.  终端输入`go test -bench . -benchmem`，查看函数每次操作需要的时间，内存，和数据交互的大小。
##### BenchMark成绩
```
BenchmarkDataHandler-12                1        16364282400 ns/op          61600 B/op       1422 allocs/op
BenchmarkCalculation-12                1        10723325200 ns/op             96 B/op          1 allocs/op
BenchmarkProcessData-12            16731             72945 ns/op               0 B/op          0 allocs/op
PASS
```
从以上结果，可以发现瓶颈出现在哪个函数。<br>
```go
package pkg

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func DataHandler(c *gin.Context) {
	// 增加一些低效计算，以增加 CPU 占用
	for i := 0; i < 1000; i++ { // 增加一个循环，增加计算量
		_ = ProcessData(i) // 调用低效计算函数
	}

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
	total := 0
	for i := 0; i < dataSize; i++ {
		// 每个循环有延时操作，模拟高延迟
		time.Sleep(time.Millisecond * 10) // 延迟 10 毫秒
		total += ProcessData(i)           // 调用低效的函数进行冗余计算
	}
	return total
}

func ProcessData(x int) int {
	// 一个低效的计算方法，用来增加计算时间
	result := 0
	// 进行大量无用的循环，增加计算耗时
	for i := 0; i < 5000; i++ { // 扩大循环次数，增加 CPU 占用
		result += rand.Intn(100)
	}
	return result + x
}
```
##### 优化后
##### Benchmark成绩
```
BenchmarkDataHandler-12               64          18880492 ns/op           88816 B/op       3102 allocs/op
BenchmarkCalculation-12               96          12152069 ns/op           56241 B/op       2003 allocs/op
BenchmarkProcessData-12            15462             85503 ns/op               0 B/op          0 allocs/op
PASS
```
```go
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
```
当要再一步优化Goroutine的执行效率时：
+ 必须避免锁的竞争(如悲观锁，`mutex` 等)
+ 使用缓存池来缓存和复用内存，避免新的goroutine都是用新的内存
+ 将任务合理拆分，确保每个goroutine都有足够大的任务进行运算，避免多于数量的goroutine创建
### 9. JSON处理与数据转换
***
为了让接收的JSON转换为Go结构体，必须使用`Go struct`，并且定义一样的JSON Key
`json:"Key"`来确保正确的格式转换。
```go
package main

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
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
```
如果created_at可能为空的话，可以对应客户需求来调整需要执行的逻辑。<br>
例如：<br>
1. 服务器自动置入目前的时间```time.Now()```。
2. 逻辑代码上传送```NULL```函数
3. 使用错误代码包(例如 ```logrus```)来反馈错误信息，并且停止或者继续此行动