package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Todo struct {
	ID       uint   `gorm:"primary_key"`
	Task     string `gorm:"not null"`
	Complete bool   `gorm:"not null;default:false"`
}

func main() {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Todo{})

	router := gin.Default()

	router.GET("/", homeHandler)
	router.GET("/todos", getTodosHandler(db))
	router.POST("/todos", createTodoHandler(db))
	router.PUT("/todos/:id", updateTodoHandler(db))
	router.PUT("/todos/:id/complete", completeTodoHandler(db))
	router.PUT("/todos/:id/uncomplete", uncompleteTodoHandler(db))
	router.DELETE("/todos/:id", deleteTodoHandler(db))
	router.GET("/alive", healthCheckHandler)
	router.GET("/metrics", metricsHandler)

	router.Run(":8080")
}

func homeHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("static/index.html")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getTodosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		c.JSON(http.StatusOK, todos)
	}
}

func createTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&todo)
		c.JSON(http.StatusCreated, todo)
	}
}

func updateTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	}
}

func completeTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		todo.Complete = true
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	}

}

func uncompleteTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		todo.Complete = false
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	}

}

func deleteTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		db.Delete(&todo)
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func metricsHandler(c *gin.Context) {
	timeStart := time.Now()

	url := c.Request.Host + c.Request.URL.Path
	status := c.Writer.Status()
	responseTime := time.Since(timeStart).Seconds() * 1000

	c.JSON(http.StatusOK, gin.H{
		"url":          url,
		"status":       status,
		"responseTime": strconv.FormatFloat(responseTime, 'f', -1, 64) + " ms",
	})
}
