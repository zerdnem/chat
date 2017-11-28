package utils

import (
	"math/rand"
	"time"
)

type Info struct {
	Nickname    string
	Color       string
	LastMessage time.Time
	ID          string
}

type Users struct {
	SocketID string
	Info     Info
}

func (u *Users) NewId(id string) string {
	u.SocketID = id
	return u.SocketID
}

func (u *Users) NewInfo(i Info) Info {
	u.Info = i
	return u.Info
}

func (u *Users) NewUser(id string, i Info) (string, Info) {
	newid := u.NewId(id)
	u.SocketID = newid
	newinfo := u.NewInfo(i)
	u.Info = newinfo
	return u.SocketID, u.Info
}

func GenerateUserColor() string {
	colors := []string{"#c0392b", "#16a085", "#27ae60", "#2980b9", "#2c3e50", "#8e44ad", "#d35400", "#f39c12", "#34495e", "#9b59b6", "#1abc9c"}
	return colors[rand.Intn(len(colors))]

}

func GenerateUserId() string {
	var chars = []rune("abcdefghjkmnopqrstvwxyz01234567890")
	id := make([]rune, 10)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)

}

func GetUserById(users []*Users) (id string) {
	var userid interface{}
	for _, v := range users {

		userid = v.SocketID
	}
	return userid.(string)

}

func RemoveUser(users []*Users, id string) (newusers []*Users) {
	for i, v := range users {
		if id == v.SocketID {
			users[i] = users[len(users)-1]
		}
	}
	return users[:len(users)-1]

}

func AddUser(nickname string, socketid string) Users {
	id := GenerateUserId()
	color := GenerateUserColor()
	userinfo := Info{
		nickname,
		color,
		time.Time{},
		id,
	}
	user := Users{}
	user.NewUser(socketid, userinfo)

	return user

}
