package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	config "codeberg.org/vnpower/pixivfe/v2/core/config"
	"github.com/tidwall/gjson"
)

type HttpResponse struct {
	Ok         bool
	StatusCode int

	Body    string
	Message string
}

func WebAPIRequest(context context.Context, URL, token string) HttpResponse {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintf("Failed to create a request to %s\n.", URL),
		}
	}
	req = req.WithContext(context)

	req.Header.Add("User-Agent", config.GlobalServerConfig.UserAgent)
	req.Header.Add("Accept-Language", config.GlobalServerConfig.AcceptLanguage)

	if token == "" {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: config.GetRandomDefaultToken(),
		})
	} else {
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintf("Failed to send a request to %s\n.", URL),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintln("Failed to parse request data."),
		}
	}

	return HttpResponse{
		Ok:         true,
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Message:    "",
	}
}

func UnwrapWebAPIRequest(context context.Context, URL, token string) (string, error) {
	resp := WebAPIRequest(context, URL, token)

	if !resp.Ok {
		return "", errors.New(resp.Message)
	}
	if !gjson.Valid(resp.Body) {
		return "", fmt.Errorf("Invalid JSON: %v", resp.Body)
	}

	err := gjson.Get(resp.Body, "error")

	if !err.Exists() {
		return "", errors.New("Incompatible request body")
	}

	if err.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	return gjson.Get(resp.Body, "body").String(), nil
}
