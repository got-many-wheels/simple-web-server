package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/goth"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUrl := "http://url-to-gravatar/"
	testChatUser := chatUser{goth.User{AvatarURL: testUrl}, ""}
	url, err := authAvatar.GetAvatarURL(testChatUser)

	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	testChatUser := chatUser{goth.User{}, "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error", err)
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "abc.jpg"), []byte{}, 0644)
	if err != nil {
		t.Fatalf("failed to write avatar file: %v", err)
	}

	fs := FileSystemAvatar{path: dir}

	testChatUser := chatUser{goth.User{}, "abc"}

	url, err := fs.GetAvatarURL(testChatUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s/abc.jpg", dir)
	if url != expected {
		t.Fatalf("unexpected url. got=%s expected=%s", url, expected)
	}
}
