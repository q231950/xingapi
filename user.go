package xingapi

import (
	"fmt"
)

// User interface defines methods to call on a user entity
type User interface {
	Name() string
	DisplayName() string
	ID() string
	ActiveEmail() string
	Birthdate() Birthdate
	BusinessAddress() Address
	fmt.Stringer
}

// XINGUser is the concrete implementation of a User
type XINGUser struct {
	InternalName            string    `json:"name"`
	InternalDisplayName     string    `json:"display_name"`
	InternalID              string    `json:"id"`
	InternalActiveEmail     string    `json:"active_email"`
	InternalBirthdate       Birthdate `json:"birth_date"`
	InternalBusinessAddress Address   `json:"business_address"`
}

func (user XINGUser) String() string {
	return " " + user.ID() + " " + user.ActiveEmail() + " - " + user.Birthdate().String() + user.BusinessAddress().String()
}

// Name return the user's full name
func (user XINGUser) Name() string {
	return user.InternalName
}

// DisplayName returns the users display name
func (user XINGUser) DisplayName() string {
	return user.InternalDisplayName
}

// ID returns the user's user ID
func (user XINGUser) ID() string {
	return user.InternalID
}

// ActiveEmail returns the user's email address
func (user XINGUser) ActiveEmail() string {
	return user.InternalActiveEmail
}

// Birthdate represents the user's birthdate
func (user XINGUser) Birthdate() Birthdate {
	return user.InternalBirthdate
}

// BusinessAddress returns the user's business address
func (user XINGUser) BusinessAddress() Address {
	return user.InternalBusinessAddress
}
