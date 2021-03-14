package gozapread

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
)

func (c *ZapClient) Balance() (uint, error) {
	if resp, err := c.client.Get(c.url + "Account/Balance"); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var resp BalanceResp
			if json.Unmarshal(body, &resp) == nil {
				if balance, err := strconv.ParseUint(resp.Balance, 10, 32); err == nil {
					return uint(balance), nil
				}
			}
		}
	}
	return 0, errors.New("Balance failed")
}
