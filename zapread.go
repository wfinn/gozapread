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

var client http.Client

func CheckClient() {
	if client.Jar == nil {
		panic("client not initialized (call Login first)")
	}
}

func Login(user, pass string) error {
	cookieJar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar:     cookieJar,
		Timeout: 30 * time.Second,
	}
	if res, err := client.Get("https://www.zapread.com/Account/Login/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			token := extractRequestVerificationToken(string(body))
			logindetails := url.Values{"__RequestVerificationToken": {token}, "UserName": {user}, "Password": {pass}, "RememberMe": {"false"}}
			if res, err := client.PostForm("https://www.zapread.com/Account/Login/", logindetails); err == nil {
				if body, err := ioutil.ReadAll(res.Body); err == nil {
					if strings.Contains(string(body), "/Account/LogOff/") { //TODO better validation?
						return nil
					}
				}
			}
		}
	}
	return errors.New("Login failed")
}

func GetGroupId(postid uint) (result uint) {
	CheckClient()
	if res, err := client.Get(fmt.Sprintf(`https://www.zapread.com/Post/Detail/%d`, postid)); err == nil {
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

func extractRequestVerificationToken(html string) (result string) {
	re := regexp.MustCompile(`<input name="__RequestVerificationToken" type="hidden" value="[^"]+`)
	if re.MatchString(html) {
		result = strings.Split(re.FindString(html), `value="`)[1]
	}
	return
}

func UnreadMessages() bool { //TODO return the uint instead
	CheckClient()
	if res, err := client.Get("https://www.zapread.com/Messages/UnreadMessages/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return !(string(body) == "0")
		}
	}
	return false
}

type MessageTable struct {
	Draw            uint `json:"draw"`
	RecordsTotal    uint `json:"recordsTotal"`
	RecordsFiltered uint `json:"recordsFiltered"`
	Data            []struct {
		ID      uint   `json:"Id"`
		Status  string `json:"Status"`
		Type    string `json:"Type"`
		From    string `json:"From"`
		FromID  string `json:"FromID"`
		Date    string `json:"Date"`
		Link    string `json:"Link"`
		Anchor  string `json:"Anchor"`
		Message string `json:"Message"`
	} `json:"data"`
}

func GetMessageTable() (MessageTable, error) {
	CheckClient()
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"From","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Message","name":"Message","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	req, err := http.NewRequest("POST", "https://www.zapread.com/Messages/GetMessagesTable", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	if err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			var messages MessageTable
			if json.Unmarshal(body, &messages) == nil {
				return messages, nil
			}
		}
	}
	return *new(MessageTable), errors.New("Blah")
}

type Post struct {
	PostID   uint   `json:"PostId"`
	Content  string `json:"Content"`
	GroupID  uint   `json:"GroupId"`
	UserID   bool   `json:"UserId"`
	Title    string `json:"Title"`
	IsDraft  bool   `json:"IsDraft"`
	Language string `json:"Language"`
}

type PostResponse struct {
	Result      string `json:"result"`
	Success     bool   `json:"success"`
	PostID      uint   `json:"postId"`
	HTMLContent string `json:"HTMLContent"`
}

func SubmitNewPost(title, content string, groupid uint) (PostResponse, error) {
	CheckClient()
	post := Post{PostID: 0, Content: content, GroupID: groupid, UserID: false, Title: title, IsDraft: false, Language: "en"}
	if j, err := json.Marshal(post); err == nil {
		if res, err := client.Get("https://www.zapread.com/Post/NewPost/"); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				token := extractRequestVerificationToken(string(body))
				if req, err := http.NewRequest("POST", "https://www.zapread.com/Post/SubmitNewPost/", bytes.NewBuffer(j)); err == nil {
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("__RequestVerificationToken", token)
					res, err := client.Do(req)
					defer res.Body.Close()
					if err == nil {
						if body, err := ioutil.ReadAll(res.Body); err == nil {
							var resp PostResponse
							if json.Unmarshal(body, &resp) == nil {
								if resp.Success {
									return resp, nil
								}
							}
						}
					}
				}
			}

		}
	}
	return *new(PostResponse), errors.New("SubmitNewPost failed")
}

func DismissMessage(id uint) error {
	CheckClient()
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if req, err := http.NewRequest("POST", "https://www.zapread.com/Messages/DismissMessage", bytes.NewBuffer([]byte(jsonStr))); err == nil {
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		defer res.Body.Close()
		if err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				if string(body) == `{"Result":"Success"}` {
					return nil
				}
			}
		}

	}
	return errors.New("DismissMessage failed")
}

type Comment struct {
	CommentContent string `json:"CommentContent"`
	PostID         uint   `json:"PostId"`
	CommentID      uint   `json:"CommentId"`
	IsReply        bool   `json:"IsReply"`
}

func AddComment(content string, postid, commentid uint) error {
	CheckClient()
	comment := Comment{CommentContent: content, PostID: postid, CommentID: commentid, IsReply: commentid != 0}
	if j, err := json.Marshal(comment); err == nil {
		if res, err := client.Get("https://www.zapread.com/?l=1"); err == nil {
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				token := extractRequestVerificationToken(string(body))
				if req, err := http.NewRequest("POST", "https://www.zapread.com/Comment/AddComment", bytes.NewBuffer(j)); err == nil {
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("__RequestVerificationToken", token)
					res, err := client.Do(req)
					defer res.Body.Close()
					if err == nil {
						if body, err := ioutil.ReadAll(res.Body); err == nil {
							if strings.Contains(string(body), `"success":true`) {
								return nil
							}
						}
					}
				}
			}

		}
	}
	return errors.New("AddComment failed")

}