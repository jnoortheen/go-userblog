package actions

import (
	"muserblog/models"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Comment)
// DB Table: Plural (Comments)
// Resource: Plural (Comments)
// Path: Plural (/comments)

// CommentsResource is the resource for the comment model
type CommentsResource struct {
	buffalo.Resource
}

// List gets all Comments. This function is mapped to the the path
// GET /comments
func (v CommentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	comments := &models.Comments{}
	// You can order your list here. Just change
	err := tx.All(comments)
	// to:
	// err := tx.Order("(case when completed then 1 else 2 end) desc, lower([sort_parameter]) asc").All(comments)
	// Don't forget to change [sort_parameter] to the parameter of
	// your model, which should be used for sorting.
	if err != nil {
		return err
	}
	return c.Render(200, r.JSON(comments))
}

// Create adds a comment to the DB. This function is mapped to the
// path POST /comments
func (v CommentsResource) Create(c buffalo.Context) error {
	// Allocate an empty Comment
	comment := &models.Comment{}
	// Bind comment to the html form elements
	err := c.Bind(comment)
	if err != nil {
		return err
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(comment)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}
	// Success!
	return c.Render(201, r.JSON(comment))
}

// Update changes a comment in the DB. This function is mapped to
// the path PUT /comments/{comment_id}
func (v CommentsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Comment
	comment := &models.Comment{}
	err := tx.Find(comment, c.Param("comment_id"))
	if err != nil {
		return err
	}
	// Bind comment to the html form elements
	err = c.Bind(comment)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(comment)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}
	// Success!
	return c.Render(200, r.JSON(comment))
}

// Destroy deletes a comment from the DB. This function is mapped
// to the path DELETE /comments/{comment_id}
func (v CommentsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Comment
	comment := &models.Comment{}
	// To find the Comment the parameter comment_id is used.
	err := tx.Find(comment, c.Param("comment_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(comment)
	if err != nil {
		return err
	}
	// Success!
	return c.Render(200, r.JSON(comment))
}
