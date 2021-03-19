/*
Package gozapread is a library for ZapRead bots.

Use gozapread.Login("username, "password") to get a ZapClient to work with.

This is far from being complete.
*/
package gozapread

import (
	"bytes"
	"encoding/json"
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

//ZapClient manages the session for a user
type ZapClient struct {
	client *http.Client
	url    string
}

//Login returns a ZapClient for user
func Login(user, pass string) (*ZapClient, error) {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 30 * time.Second,
	}
	c := &ZapClient{client: client, url: "https://www.zapread.com/"}
	token, err := c.GetNewToken()
	if err == nil {
		logindetails := url.Values{"__RequestVerificationToken": {token}, "UserName": {user}, "Password": {pass}, "RememberMe": {"false"}}
		if res, err := client.PostForm(c.url+"Account/Login/", logindetails); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				if strings.Contains(string(body), "/Account/LogOff/") { //TODO better validation?
					return c, nil
				} else {
					return &ZapClient{}, fmt.Errorf("couldn't verify a successful login")
				}
			}
		} else {
			return &ZapClient{}, fmt.Errorf("the request for logging in failed: %w", err)
		}
	}
	return &ZapClient{}, fmt.Errorf("the login failed: %w", err)
}

//GetGroupID parses the group from a post
func (c *ZapClient) GetGroupID(postid uint) uint { // return an error
	res, err := c.client.Get(fmt.Sprintf(c.url+`Post/Detail/%d`, postid))
	if err == nil {
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

//GetNewToken returns a new __RequestVerificationToken
func (c *ZapClient) GetNewToken() (string, error) {
	res, err := c.client.Get(c.url)
	if err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			re := regexp.MustCompile(`<input name="__RequestVerificationToken" type="hidden" value="[^"]+`)
			if re.MatchString(string(body)) {
				token := strings.Split(re.FindString(string(body)), `value="`)[1]
				return token, nil
			}
			return "", fmt.Errorf("GetNewToken No token found")
		}
	}
	return "", fmt.Errorf("the request to get the token failed: %w", err)
}

func (c *ZapClient) postJSON(url, jsonStr string, withcsrftoken bool) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.url+url, bytes.NewBuffer([]byte(jsonStr)))
	if err == nil {
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
		} else {
			return nil, fmt.Errorf("the request to %s failed: %w", url, err)
		}
	}
	return nil, err
}

//ParseTips gets tips from unread alerts. Hint: DismissAlert(Tip.AlertID)
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

//IsUserNameOnline uses the User/Hover endpoint to see if a user is online
func (c *ZapClient) IsUserNameOnline(name string) (bool, error) {
	user := userHover{
		UserID:   0,
		Username: name,
	}

	jsonSlc, err := json.Marshal(user)
	if err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), "Online") {
				return true, nil
			} else if strings.Contains(string(resp), "Offline") {
				return false, nil
			}

		} else {
			return false, err
		}
	}
	return false, err
}

//IsUserIDOnline uses the User/Hover endpoint to see if a user is online
func (c *ZapClient) IsUserIDOnline(id uint) (bool, error) {
	user := userHover{
		UserID:   id,
		Username: "",
	}
	//TODO refactor this, make a 3rd function containing everything below
	jsonSlc, err := json.Marshal(user)
	if err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), "Online") {
				return true, nil
			} else if strings.Contains(string(resp), "Offline") {
				return false, nil
			}

		} else {
			return false, err
		}
	}
	return false, err
}

//GetUserID uses the User/Hover endpoint to get the ID for a name
func (c *ZapClient) GetUserID(name string) (uint, error) {
	user := userHover{
		UserID:   0,
		Username: name,
	}

	jsonSlc, err := json.Marshal(user)
	if err == nil {
		if resp, err := c.postJSON("User/Hover/", string(jsonSlc), true); err == nil {
			re := regexp.MustCompile(`follow.[0-9]+`)
			if re.MatchString(string(resp)) {
				idstr := strings.Split(re.FindString(string(resp)), `follow(`)[1]
				if uid, err := strconv.ParseUint(idstr, 10, 32); err == nil {
					return uint(uid), nil
				} else {
					return 0, err
				}
			}
			return 0, fmt.Errorf("GetUserId no userid found")
		}
	}
	return 0, err
}
