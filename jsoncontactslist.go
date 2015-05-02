package xingapi

// JSONContactsList is a wrapper for parsing contacts JSON.
type JSONContactsList struct {
	JSONContactsUserIDList JSONContactsUserIDList `json:"contacts"`
}

/*
JSONContactsUserIDList is a wrapper for parsing ContactsLists that consist of
a Total number of users and a list of UserIDs
*/
type JSONContactsUserIDList struct {
	Total int                 `json:"total"`
	Users []map[string]string `json:"users"`
}

// UserIDs returns a list of UserIDs for the list's users
func (jsonList *JSONContactsUserIDList) UserIDs() []string {
	userIds := []string{}

	for _, idMap := range jsonList.Users {
		userIds = append(userIds, idMap["id"])
	}

	return userIds
}
