package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/api/auth"
	"github.com/jordanjoz/dd-vote/api/handlers"
	"github.com/jordanjoz/dd-vote/controllers"

	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"

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
	db.AutoMigrate(&table.Post{}, &table.Group{}, &table.User{}, &table.Vote{}, &table.Comment{}, &table.Permission{})

	// get api handler instances
	pc := handlers.NewPostController(db)
	gc := handlers.NewGroupController(db)
	cc := handlers.NewCommentController(db)
	uc := handlers.NewUserController(db)
	ac := handlers.NewAdminController(db)
	vc := handlers.NewVoteController(db)
	ec := handlers.NewExportController(db)

	// get view controller instances
	pvc := controllers.NewPageViewController(db)
	avc := controllers.NewAdminViewController(db)

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
		views.GET("/g/:gname", pvc.ShowGroupPage)
		views.GET("/admin/:gname", avc.ShowAdminPage)
	}

	// v1 api calls
	v1 := router.Group("api/v1")
	{
		// endpoints WITHOUT auth
		v1.POST("/login", uc.LoginWithClientID)
		v1.POST("/admin/login", ac.Login)
		v1.GET("/groups/:gname/posts", pc.GetAllPostsForGroup)

		// api v1 calls WITH auth
		v1auth := v1.Group("")
		{
			v1auth.Use(UseAuth)
			v1auth.POST("/logout", uc.Logout)
			v1auth.POST("/groups/:gname/posts", pc.CreatePost)
			v1auth.POST("/groups", gc.GetOrCreateGroup)
			v1auth.POST("/posts/:puuid/comments", cc.CreateComment)
			v1auth.POST("/posts/:puuid/votes", vc.CreateOrUpdateVote)
			v1auth.GET("/groups/:gname/votes/user", vc.GetUserVotes)
			v1auth.GET("/groups/:gname/export/all", ec.GetAllQuestionsCSV)
		}
	}

	router.Run(":8081")
}

// UseAuth rejects unauthorized api requests
func UseAuth(c *gin.Context) {
	if !auth.IsLoggedIn(c) {
		c.JSON(http.StatusUnauthorized, resp.APIResponse{IsError: false, Message: "User is not logged in"})
		c.Abort()
	}
}
