package xingapi

import (
	"fmt"
	"net/http"

	"github.com/str1ngs/ansi/color"
)

// Printer can be used to print out various domain specific objects in an individually formatted way.
type Printer struct{}

// PrintResponse prints a HTTP response
func PrintResponse(response *http.Response) {
	var colorCode string
	if response.StatusCode == 200 {
		colorCode = "g"
	} else {
		colorCode = "r"
	}
	color.Printf(colorCode, fmt.Sprintf("%s\n", response.Status))
}

// PrintCommand prints a command formatted
func PrintCommand(command string) {
	color.Printf("c", "%s", command)
}

// PrintError prints out an error
func PrintError(err error) {
	color.Printf("r", "%s", err.Error())
}

// PrintMessageWithParam allows to print two arbitrary strings
func PrintMessageWithParam(message string, param string) {
	color.Printf("m", "%s%s", message, param)
}

// Print prints an arbitrary string
func Print(s string) {
	color.Printf("d", s)
}

// PrintUser prints out a User in a verbose style
func PrintUser(user User) {
	color.Printf("", fmt.Sprintf("-----------------------------------\n%s <%s>:\n", user.DisplayName(), user.ID()))
	color.Printf("d", fmt.Sprintf("Email address:\t\t%s\nDate of birth:\t\t%s\n", user.ActiveEmail(), user.Birthdate()))
}

// PrintUserOneLine prints a User in one line
func PrintUserOneLine(user User) {
	color.Printf("m", "%s <%s>\n", user.DisplayName(), user.ID())
}
