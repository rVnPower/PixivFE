package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"codeberg.org/vnpower/pixivfe/models"
)

type PixivClient struct {
	Client *http.Client

	Cookie map[string]string
	Header map[string]string
	Lang   string
}

func (p *PixivClient) SetHeader(header map[string]string) {
	p.Header = header
}

func (p *PixivClient) AddHeader(key, value string) {
	p.Header[key] = value
}

func (p *PixivClient) SetUserAgent(value string) {
	p.AddHeader("User-Agent", value)
}

func (p *PixivClient) SetCookie(cookie map[string]string) {
	p.Cookie = cookie
}

func (p *PixivClient) AddCookie(key, value string) {
	p.Cookie[key] = value
}

func (p *PixivClient) SetSessionID(value string) {
	p.Cookie["PHPSESSID"] = value
}

func (p *PixivClient) SetLang(lang string) {
	p.Lang = lang
}

func (p *PixivClient) Request(URL string, token ...string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", URL, nil)

	// Add headers
	for k, v := range p.Header {
		req.Header.Add(k, v)
	}

	for k, v := range p.Cookie {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	if token != nil {
		req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: token[0]})
	}

	// Make a request
	resp, err := p.Client.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 200 {
		return resp, errors.New(fmt.Sprintf("Pixiv returned code: %d for request %s", resp.StatusCode, URL))
	}

	return resp, nil
}

func (p *PixivClient) TextRequest(URL string, tokens ...string) (string, error) {
	var token string

	if len(token) > 0 {
		token = tokens[0]
	}

	/// Make a request to a URL and return the response's string body
	resp, err := p.Request(URL, token)
	if err != nil {
		return "", err
	}

	// Extract the bytes from server's response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body), err
	}

	return string(body), nil
}

func (p *PixivClient) PixivRequest(URL string, tokens ...string) (json.RawMessage, error) {
	/// Make a request to a Pixiv API URL with a standard response, handle errors and return the raw JSON response
	var response models.PixivResponse
	var token string

	if len(token) > 0 {
		token = tokens[0]
	}

	body, err := p.TextRequest(URL, token)
	// body = strings.ReplaceAll(body, "i.pximg.net", configs.ProxyServer)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}
	if response.Error {
		// Pixiv returned an error
		return nil, errors.New("Pixiv responded: " + response.Message)
	}

	return response.Body, nil
}
