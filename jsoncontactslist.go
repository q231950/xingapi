// jsoncontactslist.go
package xingapi

type JSONContactsList struct {
	JSONContactsUserIdList JSONContactsUserIdList `json:"contacts"`
}

type JSONContactsUserIdList struct {
	Total int `json:"total"`
	Users []map[string]string `json:"users"`
}

func (jsonList *JSONContactsUserIdList) UserIds() []string {
	userIds := []string{}

	for _, idMap := range jsonList.Users {
		userIds = append(userIds, idMap["id"])
	}

	return userIds
}