package gozapread

import (
	"errors"
	"fmt"
)

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
		}
	}
	return errors.New("LeaveGroup failed")
}
