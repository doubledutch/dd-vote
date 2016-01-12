package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/doubledutch/dd-vote/api/auth"
	"github.com/doubledutch/dd-vote/api/handlers"
	"github.com/doubledutch/dd-vote/controllers"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/doubledutch/dd-vote/api/models/resp"
	"github.com/doubledutch/dd-vote/api/models/table"

	_ "github.com/lib/pq"
)

func main() {

	// connect to db
	db, err := gorm.Open("postgres", getPostgresConn())
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
	ph := handlers.NewPostHandler(db)
	gh := handlers.NewGroupHandler(db)
	ch := handlers.NewCommentHandler(db)
	uh := handlers.NewUserHandler(db)
	ah := handlers.NewAdminHandler(db)
	vh := handlers.NewVoteHandler(db)
	eh := handlers.NewExportHandler(db)

	// get view controller instances
	pvc := controllers.NewPageController(db)
	avc := controllers.NewAdminController(db)

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
		v1.POST("/login", uh.LoginWithClientID)
		v1.POST("/admin/login", ah.Login)
		v1.GET("/groups/:gname/posts", ph.GetAllPostsForGroup)

		// api v1 calls WITH auth
		v1auth := v1.Group("")
		{
			v1auth.Use(UseAuth)
			v1auth.POST("/logout", uh.Logout)
			v1auth.POST("/groups/:gname/posts", ph.CreatePost)
			v1auth.DELETE("/posts/:puuid", ph.DeletePost)
			v1auth.POST("/groups", gh.GetOrCreateGroup)
			v1auth.POST("/posts/:puuid/comments", ch.CreateComment)
			v1auth.POST("/posts/:puuid/votes", vh.CreateOrUpdateVote)
			v1auth.GET("/groups/:gname/votes/user", vh.GetUserVotes)
			v1auth.GET("/groups/:gname/export/all", eh.GetAllQuestionsCSV)
			v1auth.GET("/groups/:gname/export/top", eh.GetTopUsersCSV)
		}
	}

	router.Run(":8081")
}

func getPostgresConn() string {
	conn := os.Getenv("DB_CONN")
	if conn != "" {
		return conn
	}

	host := os.Getenv("POSTGRES_ADDR")
	if host == "" {
		host = os.Getenv("DDVOTE_DB_PORT_5432_TCP_ADDR")
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = os.Getenv("DDVOTE_DB_PORT_5432_TCP_PORT")
	}
	username := os.Getenv("POSTGRES_USERNAME")
	database := os.Getenv("POSTGRES_DATABASE")

	conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s ",
		host, port, username, database)

	password := os.Getenv("POSTGRES_PASSWORD")
	if password != "" {
		conn += fmt.Sprintf(" password=%s", password)
	}

	// Assume ssl is disabled for now
	conn += " sslmode=disable"
	return conn
}

// UseAuth rejects unauthorized api requests
func UseAuth(c *gin.Context) {
	if !auth.IsLoggedIn(c) {
		c.JSON(http.StatusUnauthorized, resp.APIResponse{IsError: false, Message: "User is not logged in"})
		c.Abort()
	}
}
