package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gplus"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/auth0/callback"), os.Getenv("AUTH0_DOMAIN")),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/gplus/callback")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/facebook/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}
