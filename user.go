// user.go
package xingapi

import (
	"fmt"
)

type User interface {
	Name() string
	DisplayName() string
	Id() string
	ActiveEmail() string
	Birthdate() Birthdate
	BusinessAddress() Address
	fmt.Stringer
}

type XINGUser struct {
	InternalName            string    `json:"name"`
	InternalDisplayName     string    `json:"display_name"`
	InternalId              string    `json:"id"`
	InternalActiveEmail     string    `json:"active_email"`
	InternalBirthdate       Birthdate `json:"birth_date"`
	InternalBusinessAddress Address   `json:"business_address"`
}

func (user XINGUser) String() string {
	return " " + user.Id() + " " + user.ActiveEmail() + " - " + user.Birthdate().String() + user.BusinessAddress().String()
}

func (user XINGUser) Name() string {
	return user.InternalName
}

func (user XINGUser) DisplayName() string {
	return user.InternalDisplayName
}

func (user XINGUser) Id() string {
	return user.InternalId
}

func (user XINGUser) ActiveEmail() string {
	return user.InternalActiveEmail
}

func (user XINGUser) Birthdate() Birthdate {
	return user.InternalBirthdate
}

func (user XINGUser) BusinessAddress() Address {
	return user.InternalBusinessAddress
}
