package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	// set a value
	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]any{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]any{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
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

	c := &client{
		userData: map[string]any{"userid": "abc"},
	}

	url, err := fs.GetAvatarURL(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s/abc.jpg", dir)
	if url != expected {
		t.Fatalf("unexpected url. got=%s expected=%s", url, expected)
	}
}
