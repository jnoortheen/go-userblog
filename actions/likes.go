package actions

import (
	"github.com/gobuffalo/buffalo"
	"muserblog/models"
	"net/http"
	"github.com/pkg/errors"
	"github.com/markbates/pop"
)

// LikeUpdate adds current user to the like table for the post or removes if exists already
func LikeUpdate(c buffalo.Context) error {
	// check user is logged in
	if c.Value("user") == nil {
		return c.Error(http.StatusUnauthorized, errors.New("User must be signed in to like post"))
	}
	user := c.Value("user").(*models.User)
	// create like object
	like := &models.Like{}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Post
	post := &models.Post{}
	// get post from database
	err := tx.Find(post, c.Request().PostFormValue("post_id"))
	if err != nil {
		return err
	}
	// check if the like record is there already
	if post.LikedBy(tx, user) {
		// case 1: unlike the post : delete the like record
		err = tx.BelongsTo(post).BelongsTo(user).First(like)
		if err != nil {
			return err
		}
		tx.Destroy(like)
	} else {
		// case 2: like the post : create a like record
		like.UserID = user.ID
		like.PostID = post.ID
		verrs, err := tx.ValidateAndCreate(like)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return c.Error(http.StatusUnprocessableEntity, verrs)
		}
	}
	// return number of likes to the post
	return c.Render(http.StatusOK, r.JSON(map[string]interface{}{"count": post.LikesCount(tx)}))
}
