package xingapi

// DummyUser is used as a mock in tests that need a user mock
type DummyUser struct {
	XINGUser
	UserID string
}

func (user DummyUser) String() string {
	return "dummy user"
}

func (user DummyUser) DisplayName() string {
	return user.UserID
}
