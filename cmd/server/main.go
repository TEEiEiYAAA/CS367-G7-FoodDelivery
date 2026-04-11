package main

import (
	"CS367-G7-FoodDelivery/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()

	// จุดเริ่มต้นสำหรับเพื่อนๆ มาเพิ่ม Route ของตัวเอง
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
