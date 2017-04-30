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
	// UniqUserNameErrMsg error message
	UniqUserNameErrMsg = "Username is already in use"
	// CheckPwdErrMsg error msg
	CheckPwdErrMsg = "Check username or password"
	// keep this key separately
	encryptKey = "5cdaae38582e3f1f9a17f7025"
)

// User db table struct
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
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	errors := validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Pwd, Name: "Pwd"},
	)
	return errors, nil
}

// WithName finds and fills user struct from DB using name
func (u *User) WithName(tx *pop.Connection) error {
	return tx.Where("name = ?", u.Name).First(u)
}

// CheckPassword call this method on a new struct with plain password. It looks up in DB for username and
// checks with salted password
func (u *User) CheckPassword(tx *pop.Connection) *validate.Errors {
	plainPwd := u.Pwd
	verrs := validate.NewErrors()
	if err := u.WithName(tx); err != nil {
		verrs.Add(validators.GenerateKey("Pwd"), CheckPwdErrMsg)
	} else if u.Pwd != (uuid.NewV3(u.Salt, plainPwd).String()) {
		verrs.Add(validators.GenerateKey("Pwd"), CheckPwdErrMsg)
	}

	return verrs
}

// ValidateUniqUsername validates the username is unique and already not taken
func (u *User) ValidateUniqUsername(tx *pop.Connection, errors *validate.Errors) (*validate.Errors, error) {
	// check username is unique
	if cnt, err := tx.Where("name = ?", u.Name).Count(User{}); cnt > 0 && err == nil {
		errors.Add(validators.GenerateKey("Name"), UniqUserNameErrMsg)
	}
	return errors, nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (u *User) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return u.ValidateUniqUsername(tx, validate.NewErrors())
}

// ValidateCreate gets run everytime you call "pop.ValidateCreate" method.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return u.ValidateUniqUsername(tx, validate.NewErrors())
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return u.ValidateUniqUsername(tx, validate.NewErrors())
}

// SaltPassword set value for salt field and changes the plain text password to salted password in struct
func (u *User) SaltPassword() {
	u.Salt = uuid.NewV1()
	u.Pwd = uuid.NewV3(u.Salt, u.Pwd).String()
}

// AuthToken generates a new encrypted authentication token with user id
func (u *User) AuthToken() string {
	strID := u.ID.String()
	return fmt.Sprintf("%s|%s", strID, encryptMsg([]byte(strID)))
}

// NewUser get a new user struct that can directly be used to create user record using pop
func NewUser(name string, plainPwd string, email string) *User {
	user := &User{Name: name, Email: nulls.String{String: email, Valid: true}, Pwd: plainPwd}
	user.SaltPassword()
	return user
}

// EncryptMsg will encrypt the given message with default key
func encryptMsg(msg []byte) []byte {
	mac := hmac.New(sha256.New, []byte(encryptKey))
	mac.Write([]byte(msg))
	return mac.Sum(nil)
}
