package gozapread

import (
	"fmt"
	"net/url"
)

//UpdateAboutMe implements Manage/UpdateAboutMe/
func (c *ZapClient) UpdateAboutMe(aboutme string) error {
	token, err := c.GetNewToken()
	if err == nil {
		values := url.Values{"__RequestVerificationToken": {token}, "AboutMe": {aboutme}}
		if _, err := c.client.PostForm(c.url+"Manage/UpdateAboutMe/", values); err == nil {
			return nil
		} else {
			return fmt.Errorf("the request to Manage/UpdateAboutMe failed: %w", err)
		}
	}
	return err
}
