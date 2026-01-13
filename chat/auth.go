package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	goth.User
	uniqueId string
}

func (c chatUser) UniqueID() string {
	return c.uniqueId
}

func (c chatUser) AvatarURL() string {
	return c.User.AvatarURL
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]

	// set the provider to be url query instead of literal path
	q := r.URL.Query()
	q.Set("provider", provider)
	r.URL.RawQuery = q.Encode()

	switch action {
	case "login":
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			log.Println(gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	case "callback":
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// handle name fallback
		name := strings.TrimSpace(user.Name)
		if len(name) == 0 {
			name = strings.TrimSpace(user.NickName)
		}
		if len(name) == 0 {
			name = strings.TrimSpace(user.FirstName + " " + user.LastName)
		}
		if len(name) == 0 {
			name = "<No Username>"
		}

		// hash email to be used as user id
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email))
		uniqueId := fmt.Sprintf("%x", m.Sum(nil))

		u := chatUser{user, uniqueId}
		avatarUrl, err := avatars.GetAvatarURL(u)
		if err != nil {
			log.Fatalf("Error while trying to get the user avatar url %v", err)
		}

		authCookieVal := objx.New(map[string]any{
			"userid":     u.UniqueID(),
			"name":       name,
			"email":      u.Email,
			"avatar_url": avatarUrl,
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieVal,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			w.Write([]byte("No auth cookie"))
			return
		}
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		MaxAge: -1,
	})

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
