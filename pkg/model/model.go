package model

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Err struct {
	Msg	string
	Err	string
}
