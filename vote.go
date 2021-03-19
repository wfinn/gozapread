package gozapread

import (
	"fmt"
	"strings"
)

//VotePost implements Vote/Post
func (c *ZapClient) VotePost(postid int, upvote bool, amount uint) error {
	up := 0
	if upvote {
		up = 1
	}
	jsonStr := fmt.Sprintf(`{"Id":%d,"d":%d,"a":%d,"tx":0}`, postid, up, amount)
	if resp, err := c.postJSON("Vote/Post", jsonStr, true); err == nil {
		if strings.Contains(string(resp), `"success":true`) {
			return nil
		} else {
			return fmt.Errorf("%svoting the post failed: %w", updown(upvote), err)
		}
	} else {
		return err
	}
}

//VoteComment implements Vote/Comment
func (c *ZapClient) VoteComment(commentid int, upvote bool, amount uint) error {
	up := 0
	if upvote {
		up = 1
	}
	jsonStr := fmt.Sprintf(`{"Id":%d,"d":%d,"a":%d,"tx":0}`, commentid, up, amount)
	if resp, err := c.postJSON("Vote/Comment", jsonStr, true); err == nil {
		if strings.Contains(string(resp), `"success":true`) {
			return nil
		} else {
			return fmt.Errorf("%svoting the comment failed: %w", updown(upvote), err)
		}
	} else {

		return err
	}
}

func updown(up bool) string {
	updown := "up"
	if !up {
		updown = "down"
	}
	return updown
}
