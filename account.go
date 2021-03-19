package gozapread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

//Balance implements Account/Balance
func (c *ZapClient) Balance() (uint, error) {
	if resp, err := c.client.Get(c.url + "Account/Balance"); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return 0, fmt.Errorf("couldn't read the response: %w", err)
		} else {
			var resp balanceResp
			if json.Unmarshal(body, &resp) == nil {
				if balance, err := strconv.ParseUint(resp.Balance, 10, 32); err == nil {
					return uint(balance), nil
				} else {
					return 0, err
				}
			} else {
				return 0, err
			}
		}
	} else {
		return 0, fmt.Errorf("couldn't fetch the balance: %w", err)
	}
}
