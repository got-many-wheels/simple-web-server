package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
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

	_, err = os.Stat(path.Join(wd, "avatars"))
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(wd, "avatars"), 0777); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
