package main

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/stretchr/objx"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, _, err := r.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = os.Stat(path.Join(wd, AVATARS_FOLDER_PATH))
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(wd, AVATARS_FOLDER_PATH), 0777); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filename := path.Join(AVATARS_FOLDER_PATH, userId)
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update the cookie with the new avatar_url path
	if authCookie, err := r.Cookie("auth"); err == nil {
		userData := objx.MustFromBase64(authCookie.Value)
		userData["avatar_url"] = filename
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: userData.MustBase64(),
			Path:  "/",
		})
	}

	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
