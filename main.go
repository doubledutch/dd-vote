package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/controllers"
	"github.com/jordanjoz/dd-vote/models"

	_ "github.com/lib/pq"
)

func main() {

	db, err := gorm.Open("postgres", "password=mysecretpassword host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Databas open error. Error:", err.Error())
	}
	log.Println("Database openned!")
	err = db.DB().Ping()
	if err != nil {
		log.Fatal("Database ping error! Error:", err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// migrations
	db.AutoMigrate(&models.Post{}, &models.Group{})

	router := gin.Default()
	router.Use(Logger())

	// Simple group: v1
	v1 := router.Group("api/v1")
	{
		//v1.POST("/login", loginEndpoint)
		//v1.POST("/submit", submitEndpoint)
		//v1.POST("/read", readEndpoint)
		v1.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// get controller instances
		pc := controllers.NewPostController(db)
		gc := controllers.NewGroupController(db)
		v1.POST("/post", pc.CreatePost)
		v1.POST("/group", gc.GetOrCreateGroup)

	}
	router.Run(":8080") // listen and serve on 0.0.0.0:8080

	// // not much here; it'll grow as we externalize config and add options
	// svc := service.TodoService{}
	// svc.Run()
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
