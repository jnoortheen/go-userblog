package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/i18n"

	"muserblog/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/sessions"
)

var secret = envy.Get("SESSION_SECRET", "$RT4rt@90")
var app *buffalo.App

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

// Store session store
var Store = sessions.NewCookieStore([]byte(secret))

// T language translator middleware provider
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:          ENV,
			SessionName:  "_muserblog_session",
			SessionStore: Store,
		})
		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		if ENV != "test" {
			// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
			app.Use(middleware.CSRF)
		}

		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		T, err = i18n.New(packr.NewBox("../locales"), "en-US")
		if err != nil {
			log.Fatal(err)
		}
		app.Use(T.Middleware())

		app.Use(Authorizer)
		app.Use(URLParamsToContextMw)

		app.Resource("/comments", CommentsResource{&buffalo.BaseResource{}})

		var postResource buffalo.Resource
		postResource = &PostsResource{&buffalo.BaseResource{}}
		app.GET("/", postResource.List)

		app.ServeFiles("/assets", packr.NewBox("../public/assets"))

		posts := app.Resource("/posts", postResource)
		posts.Use(PostsAuthorizer)
		posts.Middleware.Skip(PostsAuthorizer, postResource.List, postResource.Show)

		auth := app.Group("/auth")
		auth.POST("/{action}", AuthHandler)
		auth.GET("/{action}", AuthFormHandler)

		app.POST("/like", LikeUpdate)
	}
	return app
}
