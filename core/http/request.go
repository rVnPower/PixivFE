package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	core "codeberg.org/vnpower/pixivfe/v2/core/config"
	"github.com/tidwall/gjson"
)

type HttpResponse struct {
	Ok         bool
	StatusCode int

	Body    string
	Message string
}

func WebAPIRequest(URL string) HttpResponse {
	req, _ := http.NewRequest("GET", URL, nil)

	req.Header.Add("User-Agent", core.GlobalServerConfig.UserAgent)
	req.Header.Add("Accept-Language", core.GlobalServerConfig.AcceptLanguage)

	req.AddCookie(&http.Cookie{
		Name:  "PHPSESSID",
		Value: core.GlobalServerConfig.Token,
	})

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return HttpResponse{
			Ok:         false,
			StatusCode: 0,
			Body:       "",
			Message:    fmt.Sprintf("Failed to create a request to %s\n.", URL),
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

func UnwrapWebAPIRequest(URL, token string) (string, error) {
	resp := WebAPIRequest(URL)

	if !resp.Ok {
		return "", errors.New(resp.Message)
	}

	err := gjson.Get(resp.Body, "error")

	if !err.Exists() {
		return "", errors.New("Incompatible request body.")
	}

	if err.Bool() {
		return "", errors.New(gjson.Get(resp.Body, "message").String())
	}

	return gjson.Get(resp.Body, "body").String(), nil
}
