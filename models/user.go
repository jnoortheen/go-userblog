package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

const (
	UniqUserNameErrMsg = "Username is already in use"
	// keep this key separate when you are serious
	encryptKey = "5cdaae38582e3f1f9a17f7025"
)

type User struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	Name      string       `json:"name" db:"name"`
	Pwd       string       `json:"pwd" db:"pwd"`
	Email     nulls.String `json:"email" db:"email"`
	Salt      uuid.UUID    `json:"salt" db:"salt"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	errors := validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Pwd, Name: "Pwd"},
	)
	// check username is unique
	if cnt, err := tx.Where("name = ?", u.Name).Count(User{}); cnt > 0 && err == nil {
		errors.Add("Username", UniqUserNameErrMsg)
	}
	return errors, nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SaltPassword set value for salt field and changes the plain text password to salted password in struct
func (u *User) SaltPassword() {
	u.Salt = uuid.NewV1()
	u.Pwd = uuid.NewV3(u.Salt, u.Pwd).String()
}

// AuthToken generates a new authentication token every time called
func (u *User) AuthToken() string {
	strId := u.ID.String()
	return fmt.Sprintf("%s|%s", strId, encryptMsg([]byte(strId)))
}

// NewUser get a new user struct that can directly be used to create user record using pop
func NewUser(name string, plainPwd string, email string) *User {
	user := &User{Name: name, Email: nulls.String{String: email, Valid: true}, Pwd: plainPwd}
	user.SaltPassword()
	return user
}

// EncryptMsg will encrypt the given message with default key
func encryptMsg(msg []byte) []byte{
	mac := hmac.New(sha256.New, []byte(encryptKey))
	mac.Write([]byte(msg))
	return mac.Sum(nil)
}
