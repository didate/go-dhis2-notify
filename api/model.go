package api

import "github.com/didate/go-dhis2-notify/tools"

type Response struct {
	Pager Page   `json:"pager"`
	Users []User `json:"users"`
}

type Page struct {
	Page      int    `json:"page"`
	PageCount int    `json:"pageCount"`
	PageSize  int    `json:"pageSize"`
	Total     int    `json:"total"`
	NextPage  string `json:"nextPage"`
	PrevPage  string `json:"prevPage"`
}

type User struct {
	ID             string           `json:"id"`
	DisplayName    string           `json:"displayName"`
	PhoneNumber    string           `json:"phoneNumber"`
	LastUpdated    tools.CustomTime `json:"lastUpdated"`
	Created        tools.CustomTime `json:"created"`
	UserCredential UserCredendial   `json:"userCredentials"`
}

type Userinfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserCredendial struct {
	UserRoles     []Role   `json:"userRoles"`
	CreatedBy     Userinfo `json:"createdBy"`
	LastUpdatedBy Userinfo `json:"lastUpdatedBy"`
	Username      string   `json:"username"`
}

type Role struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}
