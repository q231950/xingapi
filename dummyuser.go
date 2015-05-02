package xingapi

// DummyUser is used as a mock in tests that need a user mock
type DummyUser struct {
	XINGUser
}

func (user DummyUser) String() string {
	return "dummy user"
}
