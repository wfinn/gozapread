package gozapread

import (
	"fmt"
)

//JoinGroup implements Group/JoinGroup
func (c *ZapClient) JoinGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/JoinGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		} else {
			return fmt.Errorf("joining the group wasn't successful: %w", err)
		}
	} else {
		return err
	}
}

//LeaveGroup implements Group/LeaveGroup
func (c *ZapClient) LeaveGroup(groupid uint) error {
	jsonStr := fmt.Sprintf(`{"gid":%d}`, groupid)
	if resp, err := c.postJSON("Group/LeaveGroup/", jsonStr, true); err == nil {
		if string(resp) == `{"success":true}` {
			return nil
		} else {
			return fmt.Errorf("leaving the group wasn't successful")
		}
	} else {
		return err
	}
}
