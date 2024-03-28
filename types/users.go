package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type UserParams struct {
	FirstName string ` json:"firstname"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUserParams(params UserParams) (*User, error) {
	encrp, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encrp),
	}, nil
}

func (u *UserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(u.FirstName) < minFirstNameLen {
		errors["firstName"] = (fmt.Sprintf("first Name has length lesser than %v please try again", minFirstNameLen))
	}
	if len(u.LastName) < minLastNameLen {
		errors["LastName"] = (fmt.Sprintf("Last Name has length lesser than %v please try again", minLastNameLen))
	}
	if len(u.Password) < minPasswordLen {
		errors["Password"] = (fmt.Sprintf("password has length lesser than %v please try again", minPasswordLen))
	}
	if !isValidEmail(u.Email) {
		errors["Email"] = ("email is not valid")
	}
	return errors
}

func isValidEmail(email string) bool {
	// Regular expression pattern for validating email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the email matches the pattern
	return regex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}
