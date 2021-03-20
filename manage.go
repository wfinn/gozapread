package gozapread

import (
	"fmt"
	"net/url"
)

//TipUser implements Manage/TipUser
func (c *ZapClient) TipUser(userid, amount uint) error {
	jsonStr := fmt.Sprintf(`{"id":%d,"amount":%d,"tx":null}`, userid, amount)
	if resp, err := c.postJSON("Manage/TipUser", jsonStr, true); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		} else {
			return fmt.Errorf("tipping the user wasn't successful")
		}
	} else {
		return err
	}
}

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
