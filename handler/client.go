package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (p *PixivClient) Request(URL string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", URL, nil)

	// Add headers
	for k, v := range p.Header {
		req.Header.Add(k, v)
	}
	for k, v := range p.Cookie {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
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

func (p *PixivClient) TextRequest(URL string) (string, error) {
	resp, err := p.Request(URL)
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
