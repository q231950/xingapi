package xingapi

// Address represents user and companys' addresses.
type Address struct {
	Street   string
	Zipcode  string `json:"zip_code"`
	City     string
	Province string
	Country  string
}

// String makes Address conform to Stringer interface.
func (address Address) String() string {
	return address.Street + address.Province + address.City + address.Zipcode + address.Country
}
