// address.go

package xingapi

// Address model for user and company entities
type Address struct {
	Street   string
	Zipcode  string `json:"zip_code"`
	City     string
	Province string
	Country  string
}

func (a Address) String() string {
	return a.Street + a.Province + a.City + a.Zipcode + a.Country
}
