package actions

import (
	"muserblog/models"
	"net/http"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// Following naming logic is implemented in Buffalo:
// Model: Singular (Post)
// DB Table: Plural (Posts)
// Resource: Plural (Posts)
// Path: Plural (/posts)
// View Template Folder: Plural (/templates/posts/)

// PostsResource is the resource for the post model
type PostsResource struct {
	buffalo.Resource
}

// URLParamsToContextMw set all url params to context and more
func URLParamsToContextMw(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if urlParams, ok := c.Params().(url.Values); ok {
			for param, val := range urlParams {
				c.Set(param, val[0])
				tx := c.Value("tx").(*pop.Connection)
				if param == "post_id" {
					post := &models.Post{}
					if err := tx.Find(post, c.Param("post_id")); err != nil {
						return c.Error(http.StatusNotFound, errors.New("Post not found"))
					}
					c.Set("post", post)
				} else if param == "comment_id" {
					comment := &models.Comment{}
					if err := tx.Find(comment, c.Param("comment_id")); err != nil {
						return c.Error(http.StatusNotFound, errors.New("Comment not found"))
					}
					c.Set("comment", comment)
				}
			}
		}
		return next(c)
	}
}

// List gets all Posts. This function is mapped to the the path
// GET /posts
func (v PostsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	posts := &models.Posts{}
	// You can order your list here. Just change
	err := tx.Order("updated_at desc").All(posts)
	// to:
	// err := tx.Order("(case when completed then 1 else 2 end) desc, lower([sort_parameter]) asc").All(posts)
	// Don't forget to change [sort_parameter] to the parameter of
	// your model, which should be used for sorting.
	if err != nil {
		return err
	}
	// Make posts available inside the html template
	c.Set("posts", posts)
	// Return post's author name
	c.Set("author_name", func(p models.Post, tx *pop.Connection) string {
		return (&p).Author(tx).Name
	})
	return c.Render(200, r.HTML("posts/index.html"))
}

// Show gets the data for one Post. This function is mapped to
// the path GET /posts/{post_id}
func (v PostsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Post
	post := &models.Post{}
	// To find the Post the parameter post_id is used.
	err := tx.Find(post, c.Param("post_id"))
	if err != nil {
		return err
	}

	// check for author
	user, ok := c.Value("user").(*models.User)
	c.Set("editablePost", (ok && post.UserID == user.ID))
	c.Set("pageTitle", post.Title)

	// Make post available inside the html template
	c.Set("post", post)
	return c.Render(200, r.HTML("posts/show.html"))
}

// New renders the form for creating a new post.
// This function is mapped to the path GET /posts/new
func (v PostsResource) New(c buffalo.Context) error {
	// Make post available inside the html template
	c.Set("post", &models.Post{})
	return c.Render(200, r.HTML("posts/new.html"))
}

// Create adds a post to the DB. This function is mapped to the
// path POST /posts
func (v PostsResource) Create(c buffalo.Context) error {
	// Allocate an empty Post
	post := &models.Post{}
	// Bind post to the html form elements
	err := c.Bind(post)
	if err != nil {
		return err
	}

	// setting author field
	user := c.Get("user").(*models.User)
	post.UserID = user.ID

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(post)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make post available inside the html template
		c.Set("post", post)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("posts/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Post was created successfully")
	// and redirect to the posts index page
	return c.Redirect(302, "/posts/%s", post.ID)
}

// Edit renders a edit formular for a post. This function is
// mapped to the path GET /posts/{post_id}/edit
func (v PostsResource) Edit(c buffalo.Context) error {
	post := c.Value("post").(*models.Post)

	user := c.Get("user").(*models.User)
	// checking author field
	if user.ID != post.UserID {
		return c.Error(http.StatusUnauthorized, errors.New("User is not authorized to edit this post"))
	}
	// Make post available inside the html template
	c.Set("post", post)
	return c.Render(200, r.HTML("posts/edit.html"))
}

// Update changes a post in the DB. This function is mapped to
// the path PUT /posts/{post_id}
func (v PostsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// get old post
	post := c.Value("post").(*models.Post)

	// update post to have the html form elements
	if err := c.Bind(post); err != nil {
		return err
	}

	// checking author field
	user := c.Get("user").(*models.User)
	if user.ID != post.UserID {
		return c.Error(http.StatusUnauthorized, errors.New("User is not authorized to edit this post"))
	}

	verrs, err := tx.ValidateAndUpdate(post)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make post available inside the html template
		c.Set("post", post)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("posts/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Post was updated successfully")
	// and redirect to the posts index page
	return c.Redirect(302, "/posts/%s", post.ID)
}

// Destroy deletes a post from the DB. This function is mapped
// to the path DELETE /posts/{post_id}
func (v PostsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	post := c.Value("post").(*models.Post)

	// checking author field
	user := c.Get("user").(*models.User)
	if user.ID != post.UserID {
		return c.Error(http.StatusUnauthorized, errors.New("User is not authorized to edit this post"))
	}

	if err := tx.Destroy(post); err != nil {
		return err
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Post was destroyed successfully")
	// Redirect to the posts index page
	return c.Redirect(302, "/posts")
}
