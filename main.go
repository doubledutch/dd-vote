package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/doubledutch/dd-vote/api/handlers"
	"github.com/doubledutch/dd-vote/api/models/table"
	"github.com/doubledutch/dd-vote/controllers"
	"github.com/doubledutch/dd-vote/middleware"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

// Specification contains the configuration we pull from the Environment with envconfig
type Specification struct {
	SessionCookieAuthSecret string `envconfig:"session_cookie_auth_secret" default:"insecuresecret"`
	CloudSQLUsername        string `envconfig:"cloud_sql_username" default:"ddvote"`
	CloudSQLPassword        string `envconfig:"cloud_sql_password" required:"true"`
	CloudSQLDatabase        string `envconfig:"cloud_sql_database" default:"ddvote"`
	CloudSQLHost            string `envconfig:"cloud_sql_host" default:"127.0.0.1"`
	CloudSQLPort            string `envconfig:"cloud_sql_port" default:"3306"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var s Specification
	err := envconfig.Process("ddvote", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("attempting to connect to mysql with connection spec", s)

	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		s.CloudSQLUsername, s.CloudSQLPassword, s.CloudSQLHost, s.CloudSQLPort, s.CloudSQLDatabase,
	))

	defer func() {
		err := db.Close()
		if err != nil {
			// XXX(fujin): gracefully degrade/restart when no db connection is available.
			panic(err)
		}
	}()

	if err != nil {
		log.Fatal("Unable to open database:", err.Error())
	}
	if err := db.DB().Ping(); err != nil {
		log.Fatal("Unable to ping database:", err.Error())
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// run migrations
	db.AutoMigrate(
		&table.Post{},
		&table.Group{},
		&table.User{},
		&table.Vote{},
		&table.Comment{},
		&table.Permission{},
	)

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
	store := sessions.NewCookieStore([]byte(s.SessionCookieAuthSecret))
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
			v1auth.Use(middleware.UseAuth)
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

	router.Run(":80")
}
