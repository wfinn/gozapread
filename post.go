package gozapread

import (
	"encoding/json"
	"errors"
)

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
