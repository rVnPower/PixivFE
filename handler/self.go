package handler

import (
	"net/url"
	"regexp"
)

func (p *PixivClient) FollowUser(id string) error {
	formData := url.Values{}
	formData.Add("mode", "add")
	formData.Add("type", "user")
	formData.Add("user_id", id)
	formData.Add("tag", "")
	formData.Add("restrict", "0")
	formData.Add("format", "json")

	init, err := p.GetCSRF()
	println(init)
	if err != nil {
		return err
	}

	pattern := regexp.MustCompile(`.*pixiv.context.token = "([a-z0-9]{32})"?.*`)
	quotesPattern := regexp.MustCompile(`([a-z0-9]{32})`)
	token := quotesPattern.FindString(pattern.FindString(init))
	println(token)

	_, err = p.RequestWithFormData(FollowUserURL, formData, token)
	if err != nil {
		return err
	}

	return nil
}
