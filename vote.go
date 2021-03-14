package gozapread

import (
	"errors"
	"fmt"
	"strings"
)

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

func (c *ZapClient) VoteComment(commentid int, upvote bool, amount uint) error {
	up := 0
	if upvote {
		up = 1
	}
	jsonStr := fmt.Sprintf(`{"Id":%d,"d":%d,"a":%d,"tx":0}`, commentid, up, amount)
	if resp, err := c.postJSON("Vote/Comment", jsonStr, true); err == nil {
		if strings.Contains(string(resp), `"success":true`) {
			return nil
		}
	}
	return errors.New("VoteComment failed")
}
