package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

// Comment is mapped to comment db_table/form/json
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created" db:"created_at"`
	UpdatedAt time.Time `json:"modified" db:"updated_at"`
	Content   string    `json:"content" db:"content"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (c Comment) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CommentQueryJson is mapped to jquery-comments library fields
type CommentExt struct {
	*Comment
	FullName             string `json:"fullname"`
	CreatedByCurrentUser bool `json:"created_by_current_user"`
}

// String is not required by pop and may be deleted
func (c CommentExt) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// String is not required by pop and may be deleted
func (c *CommentExt) Update(tx *pop.Connection, currentUser *User) {
	if c.UserID != uuid.Nil {
		if c.FullName == "" {
			user := &User{}
			tx.Find(user, c.UserID)
			c.FullName = user.Name
		}
		c.CreatedByCurrentUser = (currentUser.ID == c.UserID)
	}
}

// Comments is not required by pop and may be deleted
type Comments []Comment

// String is not required by pop and may be deleted
func (c Comments) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (c *Comment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Content, Name: "Content"},
		&validators.StringIsPresent{Field: c.PostID.String(), Name: "PostID"},
		&validators.StringIsPresent{Field: c.UserID.String(), Name: "PostID"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (c *Comment) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (c *Comment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
