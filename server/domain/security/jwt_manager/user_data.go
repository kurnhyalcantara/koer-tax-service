package jwtmanager

type UserData struct {
	UserID         uint64   `json:"userID"`
	Username       string   `json:"username"`
	CompanyID      uint64   `json:"companyID"`
	CompanyName    string   `json:"companyName"`
	UserType       string   `json:"userType"`
	Authorities    []string `json:"authorities"`
	GroupIDs       []uint64 `json:"groupIDs"`
	RoleIDs        []uint64 `json:"roleIDs"`
	SessionID      string   `json:"sessionID"`
	DateTime       string   `json:"dateTime"`
	TokenCreatedAt string   `json:"tokenCreatedAt"`
	IdToken        string   `json:"idToken"`
}
