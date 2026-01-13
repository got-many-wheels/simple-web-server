package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong. ErrNoAvatarURL is returned if
	// the object is unable to get a URL for the specified client.
	GetAvatarURL(u ChatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// AuthAvatar.GetAvatarURL gets avatar url straight from the auth provider itself
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return u.AvatarURL(), nil
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

// GravatarAvatar.GetAvatarURL gets avatar url by generating a unique ID for each
// profile picture. In this case we are using email and turn them into unique hash
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct {
	// Absolute path of the avatars folder
	path string
}

func UseFileSystemAvatar(path string) *FileSystemAvatar {
	return &FileSystemAvatar{
		path: path,
	}
}

// FileSystemAvatar.GetAvatarURL gets avatar url from the /avatars folder
func (fs FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	userid := u.UniqueID()
	entries, err := os.ReadDir(fs.path)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if match, _ := path.Match(userid+"*", name); match {
			return filepath.Join(fs.path, name), nil
		}
	}

	return "", ErrNoAvatarURL
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
