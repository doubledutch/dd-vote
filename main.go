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
		log.Fatal("Unable to open database:", err.Error())
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatal("Unable to ping database:", err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// run migrations
	db.AutoMigrate(&models.Post{}, &models.Group{}, &models.User{}, &models.Vote{}, &models.Comment{})

	// get controller instances
	pc := controllers.NewPostController(db)
	gc := controllers.NewGroupController(db)
	cc := controllers.NewCommentController(db)
	uc := controllers.NewUserController(db)

	// init router
	router := gin.Default()

	// session management
	store := sessions.NewCookieStore([]byte("secret")) //TODO use environment variable secret
	router.Use(sessions.Sessions("ddvote_session", store))

	// api v1 calls WITHOUT auth
	v1 := router.Group("api/v1")
	{
		v1.POST("/login", uc.LoginWithClientID)
		v1.POST("/logout", uc.Logout)
	}

	// api v1 calls WITH auth
	v1auth := router.Group("api/v1")
	{
		v1auth.Use(UseAuth)
		v1auth.GET("/post", pc.GetAllPostsForGroup)
		v1auth.POST("/post", pc.CreatePost)
		v1auth.POST("/group", gc.GetOrCreateGroup)
		v1auth.POST("/comment", cc.CreateComment)
	}

	router.Run(":8080")
}

func UseAuth(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		c.JSON(401, models.ApiResponse{IsError: false, Message: "User is not logged in"})
		c.Abort()
	}
}
