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

func (c *ZapClient) UnreadMessages() bool { //TODO return the uint instead
	if res, err := c.client.Get(c.url + "Messages/UnreadMessages/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return !(string(body) == "0")
		}
	}
	return false
}

func (c *ZapClient) GetMessageTable() (MessageTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"From","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Message","name":"Message","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJSON("Messages/GetMessagesTable", jsonStr, true); err == nil {
		var messages MessageTable
		if json.Unmarshal(resp, &messages) == nil {
			return messages, nil
		}
	}
	return *new(MessageTable), errors.New("GetMessageTable failed")
}

func (c *ZapClient) SubmitNewPost(title, content string, groupid uint) (PostResp, error) {
	post := Post{PostID: 0, Content: content, GroupID: groupid, UserID: false, Title: title, IsDraft: false, Language: "en"}
	if jsonSlc, err := json.Marshal(post); err == nil {
		if resp, err := c.postJSON("Post/SubmitNewPost/", string(jsonSlc), true); err == nil {
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

func (c *ZapClient) DismissMessage(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJSON("Messages/DismissMessage", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}
	}
	return errors.New("DismissMessage failed")
}

func (c *ZapClient) AddComment(content string, postid, commentid uint) error {
	comment := Comment{CommentContent: content, PostID: postid, CommentID: commentid, IsReply: commentid != 0}
	if jsonSlc, err := json.Marshal(comment); err == nil {
		if resp, err := c.postJSON("Comment/AddComment", string(jsonSlc), true); err == nil {
			if strings.Contains(string(resp), `"success":true`) {
				return nil
			}
		}
	}
	return errors.New("AddComment failed")
}

func (c *ZapClient) VotePost(postid int, upvote bool, amount uint) error {
	up := 0
	if upvote {
		up = 1
	}
	jsonStr := fmt.Sprintf(`{"Id":%d,"d":%d,"a":%d,"tx":0}`, postid, up, amount)
	if resp, err := c.postJSON("Vote/Post", jsonStr, true); err == nil {
		if strings.Contains(string(resp), `"success":true`) {
			return nil
		}
	}
	return errors.New("VotePost failed")
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

func (c *ZapClient) GetDepositInvoice(amount uint) (string, error) {
	jsonStr := fmt.Sprintf(`{"amount":"%d","memo":"ZapRead.com deposit","anon":"0","use":"userDeposit","useId":-1,"useAction":-1}`, amount)
	if resp, err := c.postJSON("Lightning/GetDepositInvoice/", jsonStr, false); err == nil {
		var invoice Invoice
		if json.Unmarshal(resp, &invoice) == nil {
			return invoice.Invoice, nil
		}
	}
	return "", errors.New("GetDepositInvoice failed")
}

func (c *ZapClient) SubmitPaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/SubmitPaymentRequest", jsonStr, true); err == nil {
		var payment PaymentResp
		if json.Unmarshal(resp, &payment) == nil {
			return payment.Fees, nil
		}
	} else {
		return 0, err
	}
	return 0, errors.New("SubmitPaymentRequest failed")
}

func (c *ZapClient) ValidatePaymentRequest(request string) (uint, error) {
	//do basic checks on request
	jsonStr := fmt.Sprintf(`{"request":"%s"}`, request)
	if resp, err := c.postJSON("Lightning/ValidatePaymentRequest", jsonStr, true); err == nil {
		fmt.Println(string(resp))
		var req PaymentReq
		if json.Unmarshal(resp, &req) == nil {
			if amount, err := strconv.ParseUint(req.NumSatoshis, 10, 32); err == nil {
				return uint(amount), nil
			}
		}
	}
	return 0, errors.New("ValidatePaymentRequest failed")
}

func (c *ZapClient) TipUser(userid, amount uint) error {
	jsonStr := fmt.Sprintf(`{"id":%d,"amount":%d,"tx":null}`, userid, amount)
	if resp, err := c.postJSON("Manage/TipUser", jsonStr, true); err == nil {
		if string(resp) == `{"success":true,"Result":"Success"}` {
			return nil
		}
	}
	return errors.New("TipUser failed")
}

func (c *ZapClient) JoinGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/JoinGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		}
	}
	return errors.New("JoinGroup failed")
}

func (c *ZapClient) LeaveGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/LeaveGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		} else {
			return errors.New("LeaveGroup wasn't successful.")
		}
	} else {
		return fmt.Errorf("LeaveGroup: %w", err)
	}
}

func (c *ZapClient) Balance() (uint, error) {
	if resp, err := c.client.Get(c.url + "Account/Balance"); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var resp BalanceResp
			if json.Unmarshal(body, &resp) == nil {
				if balance, err := strconv.ParseUint(resp.Balance, 10, 32); err == nil {
					return uint(balance), nil
				}
			}
		}
	}
	return 0, errors.New("UserBalance failed")
}
func (c *ZapClient) GetAlertsTable() (AlertsTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Title","name":"Title","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJSON("Messages/GetAlertsTable", jsonStr, true); err == nil {
		var alerts AlertsTable
		if json.Unmarshal(resp, &alerts) == nil {
			return alerts, nil
		}
	}
	return *new(AlertsTable), errors.New("GetAlertsTable failed")
}

func (c *ZapClient) DismissAlert(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJSON("Messages/DismissAlert", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}

	}
	return errors.New("DismissAlert failed")
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

func (c *ZapClient) GetUnreadMessages() (UnreadMessages, error) {
	if req, err := http.NewRequest(http.MethodGet, c.url+"Messages/Unread?include_content=true&include_alerts=true", nil); err == nil {
		if token, err := c.GetNewToken(); err == nil {
			req.Header.Set("__RequestVerificationToken", token)
			if resp, err := c.client.Do(req); err == nil {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					var unread UnreadMessages
					if json.Unmarshal(body, &unread) == nil {
						return unread, nil
					}
				}
			}
		}
	}
	return *new(UnreadMessages), errors.New("UnreadMessages failed")
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

//Sends a private message, doesn't return the ID yet, I didn't need it yet.
func (c *ZapClient) SendMessage(message string, toID uint) error {
	msg := ChatMessage{Content: message, ID: toID, IsChat: true}
	if jsonSlc, err := json.Marshal(msg); err == nil {
		if resp, err := c.postJSON("Messages/SendMessage", string(jsonSlc), true); err == nil {
			if strings.HasPrefix(string(resp), `{"success":true,"result":"Success",`) {
				return nil
			}
		}
	}
	return errors.New("AddComment failed")
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

func (c *ZapClient) CheckPayment(req string) (bool, error) {
	invoice := InvoiceResp{
		Invoice:   req,
		IsDeposit: true,
	}
	if jsonSlc, err := json.Marshal(invoice); err == nil {
		if resp, err := c.postJSON("Lightning/CheckPayment/", string(jsonSlc), false); err == nil {
			var check PaymentCheck
			if json.Unmarshal(resp, &check) == nil {
				if check.Success {
					return check.Result, nil
				}
			}
		}
	}
	return false, errors.New("CheckPayment failed")
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
