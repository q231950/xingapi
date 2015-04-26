// usersrequest.go

package xingapi

type UsersRequestHandler func(users []*User, err error)

type UsersRequest struct {
	UserId     string
	Limit      int
	Offset     int
	Total      int
	Completion UsersRequestHandler
}

func (ur *UsersRequest) IsFinal() bool {
	return ur.Offset+ur.Limit >= ur.Total
}
