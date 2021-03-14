package gozapread

import (
	"errors"
	"fmt"
)

//TipUser implements Manage/TipUser
func (c *ZapClient) TipUser(userid, amount uint) error {
	jsonStr := fmt.Sprintf(`{"id":%d,"amount":%d,"tx":null}`, userid, amount)
	if resp, err := c.postJSON("Manage/TipUser", jsonStr, true); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}
	}
	return errors.New("TipUser failed")
}
