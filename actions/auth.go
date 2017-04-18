package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"muserblog/models"
	"net/http"
	"github.com/markbates/pop"
)

const (
	signin = "signin"
	signup = "signup"
)

var pageTitleForAction = map[string]string{signin:"Sign-in", signup: "Sign-up"}

// for the given action (signin, signout, signup) renders the html page
func AuthFormHandler(c buffalo.Context) error {
	action := c.Param("action")
	switch action {
	case "signout":
		store := c.Session()
		store.Clear()
		c.Flash().Add("success", "Signed out successfully")
		return c.Redirect(http.StatusFound, "/auth/signin")
	}

	c.Set("pageTitle", pageTitleForAction[signin])
	c.Set("user", &models.User{})

	return c.Render(200, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
}

// for the post action of (signin, signup) create/login user
func AuthHandler(c buffalo.Context) error {
	user := &models.User{}
	action := c.Param("action")

	// load username, email, pwd fields
	err := c.Bind(user)
	if err != nil {
		return err
	}

	switch action {
	case signin:
		fmt.Println("user signin", user)
	case signup:
		tx := c.Value("tx").(*pop.Connection)
		user.SaltPassword()
		verr, err := tx.ValidateAndCreate(user)
		if err != nil {
			return err
		}
		if verr.HasAny() {
			c.Set("user", user)
			c.Set("errors", verr)
			return c.Render(http.StatusUnprocessableEntity, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
		}
		c.Session().Set("auth_token", user.AuthToken())
		return c.Redirect(http.StatusOK, "/")
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/auth/%s", action))
}
