package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Todo struct {
	ID   uint   `gorm:"primary_key"`
	Task string `gorm:"not null"`
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=your_username dbname=your_database password=your_password sslmode=disable")
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
	router.DELETE("/todos/:id", deleteTodoHandler(db))
	router.GET("/alive", healthCheckHandler)
	router.GET("/metrics", metricsHandler)

	go func() {
		fmt.Print("Enter a task (or 'quit' to exit): ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			task := scanner.Text()
			if task == "quit" {
				break
			}
			todo := Todo{Task: task}
			db.Create(&todo)
			fmt.Print("Enter a task (or 'quit' to exit): ")
		}
	}()

	router.Run(":8080")
}

func homeHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")

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
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}