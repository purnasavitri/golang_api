package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID			uint		`gorm:"column:id;primaryKey"`
	Name 		string		`gorm:"column:name"`
	Email 		string 		`gorm:"column:email"`
	Age			string		`gorm:"column:age"`
	CreatedAt	time.Time	`gorm:"column:createdAt"`
	UpdatedAt	time.Time	`gorm:"column:UpdatedAt"`
}

func (User) TableName() string {
    return "user"
}

func main() {
	dsn := "root:@tcp(localhost:3306)/openapi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
	}

	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data":users})
	})

	router.POST("/users", func(c *gin.Context){
		var newUsers User

		if err := c.ShouldBindJSON(&newUsers); err != nil {
			fmt.Println("Error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&newUsers)
		c.JSON(http.StatusCreated, gin.H{"data": newUsers})
	})
	
	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
	
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.PUT("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
	
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
	
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		db.Save(&user)
	
		c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui", "data": user})
	})

	router.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
	
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
	
		db.Delete(&user)
	
		c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
	})

	router.Run(":3000")
}