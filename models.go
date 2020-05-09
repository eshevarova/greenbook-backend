package main

type User struct {
	UserId
	Interests string `json:"interests"`
	City      string `json:"city"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	IsFriend  bool   `json:"isFriend,omitempty"`
	Name
	Login
}

type CommonUser struct {
	UserId
	Interests string `json:"interests"`
	City      string `json:"city"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	IsFriend  bool   `json:"isFriend"`
	Name
}

type Friend struct {
	Id int64 `json:"id"`
	Name
}

type Login struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserId struct {
	Id int64 `json:"id"`
}

type UserFriend struct {
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
}
