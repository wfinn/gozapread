/*
Go ZapRead.com Api Implementation

Use gozapread.Login("username, "password") to get a ZapClient to work with.

This is far from being complete.
As ZapRead is currently in beta, this can break and change quite freqently.
*/
package gozapread

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ZapClient struct {
	client *http.Client
	url    string
}

func Login(user, pass string) (*ZapClient, error) {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 30 * time.Second,
	}
	c := &ZapClient{client: client, url: "https://www.zapread.com/"}
	if token, err := c.GetNewToken(); err == nil {
		logindetails := url.Values{"__RequestVerificationToken": {token}, "UserName": {user}, "Password": {pass}, "RememberMe": {"false"}}
		if res, err := client.PostForm(c.url+"Account/Login/", logindetails); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				if strings.Contains(string(body), "/Account/LogOff/") { //TODO better validation?
					return c, nil
				}
			}
		}
	}
	return &ZapClient{}, errors.New("Login failed")
}

func (c *ZapClient) GetGroupId(postid uint) uint { // return an error
	if res, err := c.client.Get(fmt.Sprintf(c.url+`Post/Detail/%d`, postid)); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			re := regexp.MustCompile(`data-groupid=\"[^\"]*`)
			if re.MatchString(string(body)) {
				if u, err := strconv.ParseUint(strings.Split(re.FindString(string(body)), `data-groupid="`)[1], 10, 32); err == nil {
					return uint(u)
				}
			}
		}
	}
	return 0
}

func (c *ZapClient) GetNewToken() (string, error) {
	if res, err := c.client.Get(c.url); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			re := regexp.MustCompile(`<input name="__RequestVerificationToken" type="hidden" value="[^"]+`)
			if re.MatchString(string(body)) {
				token := strings.Split(re.FindString(string(body)), `value="`)[1]
				return token, nil
			}
			return "", errors.New("GetNewToken No token found")
		}
	}
	return "", errors.New("GetNewToken failed")
}

func (c *ZapClient) postJSON(url, jsonStr string, withcsrftoken bool) ([]byte, error) {
	if req, err := http.NewRequest(http.MethodPost, c.url+url, bytes.NewBuffer([]byte(jsonStr))); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if withcsrftoken {
			token, err := c.GetNewToken()
			if err != nil {
				return nil, err
			}
			req.Header.Set("__RequestVerificationToken", token)
		}
		if resp, err := c.client.Do(req); err == nil {
			defer resp.Body.Close()
			return ioutil.ReadAll(resp.Body)
		}
	}
	return nil, errors.New("postJSON failed")
}

// Parses tips from unread alerts. Hint: DismissAlert(Tip.AlertID)
func ParseTips(alerts AlertsTable) []Tip {
	var tips []Tip
	for _, alert := range alerts.Data {
		if alert.Status == "Unread" && alert.Title == "You received a tip!" {
			split := strings.Split(strings.Split(alert.Message, "/'>")[1], "</a><br/> Amount: ")
			user := split[0]
			amountStr := strings.Split(split[1], " Satoshi.")[0]
			if amount, err := strconv.ParseUint(amountStr, 10, 32); err == nil {
				tips = append(tips, Tip{From: user, Amount: uint(amount), AlertID: alert.AlertID})
			}

		}
	}
	return tips
}

func (c *ZapClient) IsUserNameOnline(name string) (bool, error) {
	user := UserHover{
		UserID:   0,
		Username: name,
	}

	if jsonSlc, err := json.Marshal(user); err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), "Online") {
				return true, nil
			} else if strings.Contains(string(resp), "Offline") {
				return false, nil
			}

		}
	}
	return false, errors.New("IsUserNameOnline failed")
}

func (c *ZapClient) IsUserIdOnline(id uint) (bool, error) {
	user := UserHover{
		UserID:   id,
		Username: "",
	}

	if jsonSlc, err := json.Marshal(user); err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), "Online") {
				return true, nil
			} else if strings.Contains(string(resp), "Offline") {
				return false, nil
			}

		}
	}
	return false, errors.New("IsUserIdOnline failed")
}

func (c *ZapClient) GetUserId(name string) (uint, error) {
	user := UserHover{
		UserID:   0,
		Username: name,
	}

	if jsonSlc, err := json.Marshal(user); err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			re := regexp.MustCompile(`follow.[0-9]+`)
			if re.MatchString(string(resp)) {
				idstr := strings.Split(re.FindString(string(resp)), `follow(`)[1]
				if uid, err := strconv.ParseUint(idstr, 10, 32); err == nil {
					return uint(uid), nil
				}
			}
			return 0, errors.New("GetUserId no userid found")
		}
	}
	return 0, errors.New("GetUserId failed")
}
