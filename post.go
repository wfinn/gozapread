package gozapread

import (
	"encoding/json"
	"fmt"
)

//SubmitNewPost implements Post/SubmitNewPost
func (c *ZapClient) SubmitNewPost(title, content string, groupid uint) (PostResp, error) {
	post := post{PostID: 0, Content: content, GroupID: groupid, UserID: false, Title: title, IsDraft: false, Language: "en"}
	jsonSlc, err := json.Marshal(post)
	if err == nil {
		if resp, err := c.postJSON("Post/SubmitNewPost/", string(jsonSlc), true); err == nil {
			var postResp PostResp
			if err := json.Unmarshal(resp, &postResp); err == nil {
				if postResp.Success {
					return postResp, nil
				} else {
					return postResp, fmt.Errorf("submitting the post failed")
				}
			} else {
				return postResp, err
			}
		} else {
			return *new(PostResp), err
		}
	}
	return *new(PostResp), err
}
