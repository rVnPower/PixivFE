package main

import (
	"net/http"
	"os"
	"pixivfe/handler"
	"testing"
	"time"
)

func generate_test_client() *handler.PixivClient {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{
		Timeout:   time.Duration(5000) * time.Millisecond,
		Transport: transport,
	}

	pc := &handler.PixivClient{
		Client: client,
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Lang:   "en",
	}

	pc.SetSessionID(os.Getenv("PIXIVFE_TOKEN"))
	pc.SetUserAgent("Mozilla/5.0")

	return pc
}

func Test_GetArtworkByID(t *testing.T) {
	client := generate_test_client()

	start := time.Now()
	_, err := client.GetArtworkByID("97422582")
	end := time.Since(start).Seconds()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Took %.2f seconds", end)
}

func Test_GetArtworkImages(t *testing.T) {
	client := generate_test_client()

	start := time.Now()
	_, err := client.GetArtworkImages("97422582")
	end := time.Since(start).Seconds()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Took %.2f seconds", end)
}

func Test_GetArtworkComments(t *testing.T) {
	client := generate_test_client()

	start := time.Now()
	_, err := client.GetArtworkComments("97422582")
	end := time.Since(start).Seconds()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Took %.2f seconds", end)
}

func Test_GetRelatedArtworks(t *testing.T) {
	client := generate_test_client()

	start := time.Now()
	_, err := client.GetRelatedArtworks("97422582")
	end := time.Since(start).Seconds()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Took %.2f seconds", end)
}

// func Test_GetArtworkGoroutine(t *testing.T) {
// }
