// birthdate.go

package xingapi 

import "fmt"

type Birthdate struct {
	Day uint
	Month uint
	Year uint
}

func (birthdate Birthdate)String() string {
	return fmt.Sprintf("%d.%d.%d", birthdate.Day, birthdate.Month, birthdate.Year)
}