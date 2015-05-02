package xingapi

/*
UsersRequestHandler represents the function to execute once the users
request has completed.
*/
type UsersRequestHandler func(users []*User, err error)

// UsersRequest represents the request to fetch users from an offset
type UsersRequest struct {
	UserID     string
	Limit      int
	Offset     int
	Total      int
	Completion UsersRequestHandler
}

// IsFinal defines if the given request will fetch the last available user.
func (request *UsersRequest) IsFinal() bool {
	return request.Offset+request.Limit >= request.Total
}
