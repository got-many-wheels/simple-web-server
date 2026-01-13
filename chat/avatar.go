package main

import (
	"errors"
	"fmt"
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
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// AuthAvatar.GetAvatarURL gets avatar url straight from the auth provider itself
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

// GravatarAvatar.GetAvatarURL gets avatar url by generating a unique ID for each
// profile picture. In this case we are using email and turn them into unique hash
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
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
func (fs FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		useridStr, ok := userid.(string)
		if !ok {
			return "", fmt.Errorf("type assertion error. got=%T expected=string", userid)
		}
		entries, err := os.ReadDir(fs.path)
		if err != nil {
			return "", err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if match, _ := path.Match(useridStr+"*", name); match {
				return filepath.Join(fs.path, name), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
