package gozapread

import (
	"errors"
	"fmt"
)

//JoinGroup implements Group/JoinGroup
func (c *ZapClient) JoinGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/JoinGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		}
	}
	return errors.New("JoinGroup failed")
}

//LeaveGroup implements Group/LeaveGroup
func (c *ZapClient) LeaveGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/LeaveGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		}
	}
	return errors.New("LeaveGroup failed")
}
