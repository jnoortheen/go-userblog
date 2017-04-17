package actions

import (
	"github.com/gobuffalo/buffalo"
	"fmt"
)

// for the given action (signin, signout, signup) renders the html page
func AuthFormHandler(c buffalo.Context) error {
	return c.Render(200, r.String(fmt.Sprintf("page#%s", c.Param("action"))))
}

// for the post action of (signin, signup) create/login user
func AuthHandler(c buffalo.Context) error {
	return nil
}