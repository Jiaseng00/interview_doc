package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func main() {
	var db *gorm.DB
	r := gin.Default()

	repo := &dbRepo{db: db}

	r.Group("/api/orders")
	{
		r.GET("", repo.GetOrder)
		r.POST("/", repo.CreateOrder)
		r.PUT("/:id", repo.UpdateOrder)
		r.DELETE("/:id", repo.DeleteOrder)
	}

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

// JSON 回复的helper
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

type Order struct {
	Id         int `gorm:"id" json:"id"`
	CustomerId int `gorm:"customer_id" json:"customer_id"`
	Status     int `gorm:"status" json:"status"`
	CreatedAt  int `gorm:"created_at" json:"created_at"`
	UpdatedAt  int `gorm:"updated_at" json:"updated_at"`
}

type dbRepo struct {
	db *gorm.DB
}

func PaginationParams(c *gin.Context) (int, int, error) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Println("limit int格式转换失败")
		return 0, 0, err
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println("offset int格式转换失败")
		return 0, 0, err
	}
	return limit, offset, nil
}

func (r *dbRepo) GetOrder(c *gin.Context) {
	var order []Order
	limit, offset, err := PaginationParams(c)
	err = r.db.Table("orders").Limit(limit).Offset(offset).Find(&order).Error
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取订单失败")
		return
	}
	JSONResponse(c, http.StatusOK, "成功获取订单", order)
}

func (r *dbRepo) CreateOrder(c *gin.Context) {
	var order Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "无效的JSON格式")
		return
	}
	err = r.db.Table("orders").Create(&order).Error
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建订单失败")
		return
	}
	JSONResponse(c, http.StatusCreated, "成功创建订单", order)
}

func (r *dbRepo) UpdateOrder(c *gin.Context) {
	var order Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "无效的JSON格式")
		return
	}
	orderIdStr := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderIdStr)
	if err != nil {
		log.Println("转换int格式失败")
		return
	}
	err = r.db.Table("orders").Where("id = ?", orderIdInt).Updates(&order).Error
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新订单失败")
		return
	}
	JSONResponse(c, http.StatusOK, "成功更新订单", nil)
}

func (r *dbRepo) DeleteOrder(c *gin.Context) {
	var order Order
	orderIdStr := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderIdStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "id int格式转换失败")
		return
	}
	err = r.db.Table("orders").Where("id = ?", orderIdInt).Delete(&order).Error // hard delete
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "删除订单失败")
		return
	}
	JSONResponse(c, http.StatusNoContent, "成功删除订单", nil)
}
