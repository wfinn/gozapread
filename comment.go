package gozapread

import (
	"encoding/json"
	"errors"
	"strings"
)

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
