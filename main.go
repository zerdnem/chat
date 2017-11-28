package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/zerdnem/chat/utils"
)

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func main() {

	users := []*utils.Users{}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {

		so.On("enter", func(data map[string]string) {

			so.Join("chat")

			//If nickname is more than 35 characters long emit
			//accepted false
			if len(users) == 0 {
				fakeuser := utils.AddUser("admin", "socketid")
				users = append(users, &fakeuser)
			}

			if len(data["nickname"]) <= 35 {
				socketid := so.Id()
				nickname := data["nickname"]
				newuser := utils.AddUser(nickname, socketid)

				var check bool

				for _, v := range users {
					if newuser.Info.Nickname == v.Info.Nickname {
						check = true
					}
				}

				// for _, v := range users {
				//If nickname is already taken emit accepted false
				//with message nickname is already taken.
				if check {
					accepted := make(map[string]interface{})
					accepted["accepted"] = false
					accepted["message"] = "The nickname is already taken."
					var message interface{} = accepted
					so.Emit("enter response", message)
				} else {

					//If nickname is accepted broadcast message nickname
					// has joined
					users = append(users, &newuser)
					accepted := make(map[string]bool)
					accepted["accepted"] = true

					broadcast := make(map[string]string)
					broadcast["message"] = newuser.Info.Nickname + " has joined."
					var message interface{} = accepted
					var broadcastmessage interface{} = broadcast

					so.BroadcastTo("chat", "info", broadcastmessage)

					so.Emit("enter response", message)
				}
				// }
			} else {
				accepted := make(map[string]interface{})
				accepted["accepted"] = false
				accepted["message"] = "The nickname is invalid (max. 35 characters)."
				var message interface{} = accepted
				so.Emit("enter response", message)
			}

		})

		so.On("message", func(data map[string]string) {
			userID := so.Id()

			var nickname string
			var color string
			var timestamp time.Time

			for _, v := range users {
				if userID == v.SocketID {
					nickname = v.Info.Nickname
					color = v.Info.Color
					timestamp = v.Info.LastMessage
				}
			}

			//Limit message to 400 characters length
			//if over said length emit message response
			//message is too long
			if len(data["message"]) <= 400 {

				messageObject := make(map[string]interface{})
				messageObject["message"] = data["message"]
				messageObject["nickname"] = nickname
				messageObject["color"] = color
				messageObject["external"] = true

				diff := time.Now().Sub(timestamp)
				d2 := time.Duration(1500) * time.Millisecond

				if diff < d2 {
					accepted := make(map[string]interface{})
					accepted["accepted"] = false
					accepted["message"] = "Slow down! You're posting too fast."

					var message interface{} = accepted
					so.Emit("message response", message)
				} else {
					for _, v := range users {
						if userID == v.SocketID {
							v.Info.LastMessage = time.Now()
						}
					}
					so.BroadcastTo("chat", "message", messageObject)
					log.Println("Broadcast => ", messageObject)

					messageObject["external"] = false
					so.Emit("message", messageObject)
				}
			} else {
				accepted := make(map[string]interface{})
				accepted["message"] = "The message is too long."
				var message interface{} = accepted
				so.Emit("message response", message)
			}

		})

		so.On("disconnection", func() {
			log.Println("on disconnect")

			userID := so.Id()
			//if user disconnects remove user from users array
			//using their SocketID then broadcast user has left
			for _, v := range users {
				if userID == v.SocketID {
					newusers := utils.RemoveUser(users, v.SocketID)

					users = newusers
					accepted := make(map[string]interface{})
					accepted["message"] = v.Info.Nickname + " has left."
					var message interface{} = accepted
					so.BroadcastTo("chat", "info", message)
				}
			}
		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)

	})

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	LoadTemplates("templates/*.html")

	http.Handle("/", r)
	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
