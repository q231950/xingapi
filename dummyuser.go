// dummyuser.go

package xingapi

type DummyUser struct {
	XINGUser
}

func (user DummyUser) String() string {
	return "dummy user"
}
