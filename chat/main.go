package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"simple-web-server/trace"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/stretchr/objx"
)

type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room where the client is chatting in
	room *room
	// userData holds informnation about the user
	userData map[string]any
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		msg.When = time.Now()
		msg.Email = c.userData["email"].(string)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]any{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load .env file")
	}
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse()

	googleAuthClientID := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	googleAuthClientSecret := os.Getenv("GOOGLE_AUTH_CLIENT_SECRET")

	goth.UseProviders(
		google.New(googleAuthClientID, googleAuthClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout) // allow debugging

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", clearCookie)
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r) // websocket endpoint

	// get the room running
	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
