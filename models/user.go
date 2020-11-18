package models

var Users map[string]User

type User struct {
	UserId string
	Username string
	Password string
}


func init() {
	Users = make(map[string]User)

}


