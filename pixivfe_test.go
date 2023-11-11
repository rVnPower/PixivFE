package main

import (
	"net/http"
	"testing"
)

func Benchmark_Main(t *testing.B) {
	r, err := http.Get("http://localhost:8282/")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_Ranking(t *testing.B) {
	r, err := http.Get("http://localhost:8282/ranking")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_Ranking_Complex(t *testing.B) {
	r, err := http.Get("http://localhost:8282/ranking?content=all&mode=daily_r18&page=1&date=20230826")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_Artwork(t *testing.B) {
	r, err := http.Get("http://localhost:8282/artworks/111157207")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_Artwork_R18(t *testing.B) {
	r, err := http.Get("http://localhost:8282/artworks/111130033")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_User_NoSocial(t *testing.B) {
	r, err := http.Get("http://localhost:8282/users/1035047")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_User_WithSocial(t *testing.B) {
	r, err := http.Get("http://localhost:8282/users/59336265")
	if r.StatusCode != 200 {
		t.Errorf("Status code not 200: was %d", r.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}
