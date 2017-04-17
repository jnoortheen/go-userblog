package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
)

// for the given action (signin, signout, signup) renders the html page
func AuthFormHandler(c buffalo.Context) error {
	action := c.Param("action")
	var pageTitle string
	switch action {
	case "signin":
		pageTitle = "Sign-in"
	case "signup":
		pageTitle = "Sign-up"
	}
	c.Set("pageTitle", pageTitle)
	return c.Render(200, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
}

// for the post action of (signin, signup) create/login user
func AuthHandler(c buffalo.Context) error {
	return nil
}
