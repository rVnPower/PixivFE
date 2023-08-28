package main

import (
	"net/http"
	"testing"
)

func Test_Main(t *testing.T) {
	_, err := http.Get("http://localhost:8282/")
	if err != nil {
		t.Error(err)
	}
	t.Logf("/ passed")
}

func Test_Settings(t *testing.T) {
	_, err := http.Get("http://localhost:8282/settings")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Settings route passed")
}

func Test_Ranking(t *testing.T) {
	_, err := http.Get("http://localhost:8282/ranking")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Ranking route passed")
}

func Test_Ranking_Complex(t *testing.T) {
	_, err := http.Get("https://pixivfe.exozy.me/ranking?content=all&mode=daily_r18&page=1&date=20230826")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Complex ranking route passed")
}

func Test_Artwork(t *testing.T) {
	_, err := http.Get("http://localhost:8282/artworks/111157207")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Artwork route passed")
}

func Test_Artwork_R18(t *testing.T) {
	_, err := http.Get("http://localhost:8282/artworks/111130033")
	if err != nil {
		t.Error(err)
	}
	t.Logf("R18 Artwork route passed")
}

func Test_User(t *testing.T) {
	_, err := http.Get("http://localhost:8282/users/1960050")
	if err != nil {
		t.Error(err)
	}
	t.Logf("User route passed")
}
