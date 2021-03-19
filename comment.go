package gozapread

import (
	"encoding/json"
	"fmt"
	"strings"
)

//AddComment implements Comment/AddComment
func (c *ZapClient) AddComment(content string, postid, commentid uint) error {
	comment := comment{CommentContent: content, PostID: postid, CommentID: commentid, IsReply: commentid != 0}
	if jsonSlc, err := json.Marshal(comment); err == nil {
		if resp, err := c.postJSON("Comment/AddComment", string(jsonSlc), true); err == nil {
			respStr := string(resp)
			if strings.Contains(respStr, `"success":true`) {
				return nil
			} else {
				return fmt.Errorf("adding the comment wasn't successful: %s %w", respStr, err)
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
