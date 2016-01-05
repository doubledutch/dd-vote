package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/controllers"
	"github.com/jordanjoz/dd-vote/models"
	"github.com/jordanjoz/dd-vote/viewcontrollers"

	_ "github.com/lib/pq"
)

func main() {

	// connect to db
	db, err := gorm.Open("postgres", "password=mysecretpassword host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Unable to open database:", err.Error())
	}
	if err := db.DB().Ping(); err != nil {
		log.Fatal("Unable to ping database:", err.Error())
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// run migrations
	db.AutoMigrate(&models.Post{}, &models.Group{}, &models.User{}, &models.Vote{}, &models.Comment{})

	// get api controller instances
	pc := controllers.NewPostController(db)
	gc := controllers.NewGroupController(db)
	cc := controllers.NewCommentController(db)
	uc := controllers.NewUserController(db)
	ac := controllers.NewAdminController(db)
	vc := controllers.NewVoteController(db)

	// get view controller instances
	pvc := controllers.NewPageViewController(db)
	avc := viewcontrollers.NewAdminViewController(db)

	// init router
	router := gin.Default()

	// serve static files
	router.Static("/css", "./static/css")
	router.Static("/js", "./static/js")
	router.Static("/img", "./static/img")

	// session management
	store := sessions.NewCookieStore([]byte("secret")) //TODO use environment variable secret
	router.Use(sessions.Sessions("ddvote_session", store))

	// view routes
	views := router.Group("")
	{
		views.GET("/g/:gid", pvc.ShowGroupPage)
		views.GET("/admin/:gid", avc.ShowAdminPage)
	}

	// v1 api calls
	v1 := router.Group("api/v1")
	{
		// endpoints WITHOUT auth
		v1.POST("/login", uc.LoginWithClientID)
		v1.POST("/admin/login", ac.Login)
		v1.GET("/post", pc.GetAllPostsForGroup)

		// api v1 calls WITH auth
		v1auth := v1.Group("")
		{
			v1auth.Use(UseAuth)
			v1auth.POST("/logout", uc.Logout)
			v1auth.POST("/post", pc.CreatePost)
			v1auth.POST("/group", gc.GetOrCreateGroup)
			v1auth.POST("/comment", cc.CreateComment)
			v1auth.POST("/vote", vc.CreateOrUpdateVote)
			v1auth.GET("/user/votes", vc.GetUserVotes)
		}
	}

	router.Run(":8080")
}

func UseAuth(c *gin.Context) {
	session := sessions.Default(c)
	// verify that user id is set
	v := session.Get("uid") // TODO - use UUID?
	if v == nil {
		c.JSON(401, models.ApiResponse{IsError: false, Message: "User is not logged in"})
		c.Abort()
	}

	//TODO -verify that user id exists?
}
