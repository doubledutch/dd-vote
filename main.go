package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
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
	db.AutoMigrate(&models.Post{}, &models.Group{}, &models.User{}, &models.Vote{}, &models.Comment{})

	router := gin.Default()

	// session management
	store := sessions.NewCookieStore([]byte("secret")) //TODO use environment variable secret
	router.Use(sessions.Sessions("ddvote_session", store))

	authorized := router.Group("/")

	authorized.Use(AuthRequired())
	{
		authorized.GET("/incr", func(c *gin.Context) {
			session := sessions.Default(c)
			var count int
			v := session.Get("count")
			if v == nil {
				count = 0
			} else {
				count = v.(int)
				count++
			}
			session.Set("count", count)
			session.Save()
			c.JSON(200, gin.H{"count": count})
		})
	}

	// Simple group: v1
	v1 := router.Group("api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// get controller instances
		pc := controllers.NewPostController(db)
		gc := controllers.NewGroupController(db)
		cc := controllers.NewCommentController(db)

		v1.GET("/post", pc.GetAllPostsForGroup)
		v1.POST("/post", pc.CreatePost)
		v1.POST("/group", gc.GetOrCreateGroup)
		v1.POST("/comment", cc.CreateComment)

	}
	router.Run(":8080")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("uid")
		if v == nil {
			c.String(200, "user is not logged in")
			return
		}
	}
}
