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

type zapclient struct {
	client http.Client
	url    string
}

func Login(user, pass string) (zapclient, error) {
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:     cookieJar,
		Timeout: 30 * time.Second,
	}
	c := zapclient{client: client, url: "https://www.zapread.com/"}
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
	return zapclient{}, errors.New("Login failed")
}

func (c zapclient) GetGroupId(postid uint) (result uint) { // return an error
	if res, err := c.client.Get(fmt.Sprintf(c.url+`Post/Detail/%d`, postid)); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			re := regexp.MustCompile(`data-groupid=\"[^\"]*`)
			if re.MatchString(string(body)) {
				if u, err := strconv.ParseUint(strings.Split(re.FindString(string(body)), `data-groupid="`)[1], 10, 32); err == nil {
					result = uint(u)
				}
			}

		}
	}
	return
}

func (c zapclient) UnreadMessages() bool { //TODO return the uint instead
	if res, err := c.client.Get(c.url + "Messages/UnreadMessages/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return !(string(body) == "0")
		}
	}
	return false
}

func (c zapclient) GetMessageTable() (MessageTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"From","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Message","name":"Message","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJson("Messages/GetMessagesTable", jsonStr, false); err == nil {
		var messages MessageTable
		if json.Unmarshal(resp, &messages) == nil {
			return messages, nil
		}
	}
	return *new(MessageTable), errors.New("GetMessageTable failed")
}

func (c zapclient) SubmitNewPost(title, content string, groupid uint) (PostResp, error) {
	post := Post{PostID: 0, Content: content, GroupID: groupid, UserID: false, Title: title, IsDraft: false, Language: "en"}
	if jsonSlc, err := json.Marshal(post); err == nil {
		if resp, err := c.postJson("Post/SubmitNewPost/", string(jsonSlc), true); err == nil {
			var postResp PostResp
			if json.Unmarshal(resp, &postResp) == nil {
				if postResp.Success {
					return postResp, nil
				}
			}
		}
	}
	return *new(PostResp), errors.New("SubmitNewPost failed")
}

func (c zapclient) DismissMessage(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJson("Messages/DismissMessage", jsonStr, false); err == nil {
				if string(resp) == `{"Result":"Success"}` {
					return nil
				}
	}
	return errors.New("DismissMessage failed")
}

func (c zapclient) AddComment(content string, postid, commentid uint) error {
	comment := Comment{CommentContent: content, PostID: postid, CommentID: commentid, IsReply: commentid != 0}
	if jsonSlc, err := json.Marshal(comment); err == nil {
		if resp, err := c.postJson("Comment/AddComment", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), `"success":true`) {
				return nil
			}
		}
	}
	return errors.New("AddComment failed")
}

func (c zapclient) VotePost(postid int, upvote bool, amount uint) error {
	up := 0
	if upvote {
		up = 1
	}
	jsonStr := fmt.Sprintf(`{"Id":%d,"d":%d,"a":%d,"tx":0}`, postid, up, amount)
	if resp, err := c.postJson("Vote/Post", jsonStr, true); err == nil {
		if strings.Contains(string(resp), `"success":true`) {
			return nil
		}
	}
	return errors.New("VotePost failed")

}

func (c zapclient) GetNewToken() (string, error) {
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

func (c zapclient) GetDepositInvoice(amount uint) (string, error) {
	jsonStr := fmt.Sprintf(`{"amount":"%d","memo":"ZapRead.com deposit","anon":"0","use":"userDeposit","useId":-1,"useAction":-1}`, amount)
	if resp, err := c.postJson("Lightning/GetDepositInvoice/", jsonStr, false); err == nil {
		//TODO Parse Invoice
		return string(resp), nil

	}
	return "", errors.New("GetDepositInvoice failed")
}

func (c zapclient) TipUser(userid, amount uint) error {
	jsonStr := fmt.Sprintf(`{"id":%d,"amount":%d,"tx":null}`, userid, amount)
	if resp, err := c.postJson("Manage/TipUser", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}
	}
	return errors.New("TipUser failed")
}

func (c zapclient) JoinGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJson("Group/JoinGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		}
	}
	return errors.New("JoinGroup failed")
}

func (c zapclient) LeaveGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJson("Group/LeaveGroup/", jsonStr, true); err != nil {
		if string(resp) == `{"success":true}` {
			return nil
		}
	}
	return errors.New("LeaveGroup failed")
}

func (c zapclient) UserBalance() (uint, error) {
	if resp, err := c.client.Get(c.url + "Account/UserBalance"); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var resp BalanceResp
			if json.Unmarshal(body, &resp) == nil {
				return resp.Balance, nil
			}
		}

	}
	return 0, errors.New("GetUserBalance failed")
}
func (c zapclient) GetAlertsTable() (AlertsTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Title","name":"Title","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJson("Messages/GetAlertsTable", jsonStr, false); err == nil {
		var alerts AlertsTable
		if json.Unmarshal(resp, &alerts) == nil {
			return alerts, nil
		}
	}
	return *new(AlertsTable), errors.New("GetAlertsTable failed")
}

func (c zapclient) DismissAlert(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJson("Messages/DismissAlert", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}

	}
	return errors.New("DismissAlert failed")
}

func (c zapclient) postJson(url, jsonStr string, withcsrftoken bool) ([]byte, error) {
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
	return nil, errors.New("HttpPost failed")
}
