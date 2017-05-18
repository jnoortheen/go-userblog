package actions

import (
	"errors"
	"fmt"
	"muserblog/models"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
)

const (
	signin           = "signin"
	signup           = "signup"
	authTokenKeyName = "auth_token"
)

var pageTitleForAction = map[string]string{signin: "Sign-in", signup: "Sign-up"}

// AuthFormHandler for the given action (signin, signout, signup) renders the html page
func AuthFormHandler(c buffalo.Context) error {
	action := c.Param("action")
	switch action {
	case "signout":
		store := c.Session()
		store.Delete(authTokenKeyName)
		store.Save()
		c.Flash().Add("success", "Signed out successfully")
		return c.Redirect(http.StatusFound, "/auth/signin")
	}

	c.Set("pageTitle", pageTitleForAction[signin])
	c.Set("user", &models.User{})

	return c.Render(200, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
}

// AuthHandler for the post action of (signin, signup) create/login user
func AuthHandler(c buffalo.Context) error {
	user := &models.User{}
	action := c.Param("action")

	// load username, email, pwd fields
	err := c.Bind(user)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	switch action {
	case signin:
		// check whether user have entered empty username/pwd into forms fields
		verr, err := user.Validate(tx)
		if err != nil {
			return err
		}
		if verr.HasAny() {
			c.Set("user", user)
			c.Set("errors", verr)
			return c.Render(http.StatusUnprocessableEntity, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
		}
		// check password
		verr = user.CheckPassword(tx)
		if verr.HasAny() {
			c.Set("user", user)
			c.Set("errors", verr)
			return c.Render(http.StatusUnprocessableEntity, r.HTML(fmt.Sprintf("auth/%s.html", c.Param("action"))))
		}
	case signup:
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
		c.Flash().Add("success", "Account created successfully!")
		return c.Redirect(http.StatusFound, "/auth/signin")
	default:
		return c.Error(http.StatusNotFound, errors.New("not found"))
	}

	if c.Request().Form.Get("rememberMe") == "true" {
		// set expiration date as 7 days for the underlying cookie implementation
		Store.MaxAge(7 * 24 * 60 * 60)
	} else {
		// set session cookie
		Store.MaxAge(0)
	}

	c.Session().Set(authTokenKeyName, user.AuthToken())
	c.Session().Save()

	return c.Redirect(http.StatusFound, "/posts")
}

// Authorizer middleware
func Authorizer(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if authID := c.Session().Get(authTokenKeyName); authID != nil {
			tx := c.Value("tx").(*pop.Connection)
			user := &models.User{}
			userID := strings.Split(authID.(string), "|")[0]
			err := tx.Find(user, userID)
			if err == nil {
				c.Set("user", user)
			} else {
				c.Session().Clear()
			}
		}
		if c.Value("user") != nil {
			c.Set("userSignedIn", true)
		} else {
			c.Set("userSignedIn", false)
		}
		return next(c)
	}
}

// PostsAuthorizer middleware redirects if user is not logged in
func PostsAuthorizer(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if c.Value("userSignedIn").(bool) {
			return c.Redirect(http.StatusFound, "/auth/signin")
		}
		return next(c)
	}
}
