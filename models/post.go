package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"strings"
	"html/template"
)

// Post db table struct
type Post struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (p Post) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Posts is not required by pop and may be deleted
type Posts []Post

// String is not required by pop and may be deleted
func (p Posts) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (p *Post) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Content, Name: "Content"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (p *Post) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (p *Post) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Author return the author of the post
func (p *Post) Author(tx *pop.Connection) *User {
	user := &User{}
	tx.Find(user, p.UserID)
	return user
}

// LikedBy return whether the given user has liked the post or not
func (p *Post) LikedBy(tx *pop.Connection, user *User) bool {
	if user != nil {
		cnt, err := tx.BelongsTo(p).BelongsTo(user).Count(&Like{})
		return err == nil && cnt > 0
	}
	return false
}

// LikesCount return number of users liked posts so far
func (p *Post) LikesCount(tx *pop.Connection) int {
	cnt, err := tx.BelongsTo(p).Count(&Like{})
	if err == nil {
		return cnt
	}
	return 0
}

// ShortContent return few lines of the post
func (p Post) ShortContent() template.HTML {
	cont := strings.Split(p.Content, "\n")
	if len(cont) > 3 {
		cont = cont[:3]
	}
	return template.HTML(strings.Join(cont, "<br/>"))
}

// ShortContent return few lines of the post
func (p *Post) ContentHtml() template.HTML {
	return template.HTML(strings.Replace(p.Content, "\n", "<br/>", -1))
}
