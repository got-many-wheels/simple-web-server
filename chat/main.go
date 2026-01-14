package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"log"
	"net/http"
	"os"
	"simple-web-server/trace"
)

const AVATARS_FOLDER_PATH = "avatars"

var avatars = TryAvatars{
	UseFileSystemAvatar(AVATARS_FOLDER_PATH),
	UseAuthAvatar,
	UseGravatarAvatar,
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
		google.New(googleAuthClientID, googleAuthClientSecret, "http://localhost:8080/auth/callback/google", "email", "profile"),
	)

	r := newRoom()

	r.tracer = trace.New(os.Stdout) // allow debugging

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", clearCookie)
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r) // websocket endpoint

	// endpoints that are responsible for user avatar upload
	http.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	http.HandleFunc("/uploader", uploadHandler)

	// static assets
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	// get the room running
	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
