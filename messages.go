package gozapread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//UnreadMessages implements the Messages/UnreadMessages. Note: currently it only returns a bool, not the number of messages
func (c *ZapClient) UnreadMessages() bool { //TODO return the uint instead
	if res, err := c.client.Get(c.url + "Messages/UnreadMessages/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return !(string(body) == "0")
		}
	}
	return false
}

//GetMessageTable implements Messages/GetMessageTable
func (c *ZapClient) GetMessageTable() (MessageTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"From","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Message","name":"Message","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	resp, err := c.postJSON("Messages/GetMessagesTable", jsonStr, true)
	if err == nil {
		var messages MessageTable
		if json.Unmarshal(resp, &messages) == nil {
			return messages, nil
		} else {
			return *new(MessageTable), err
		}
	}
	return *new(MessageTable), err
}

//DismissMessage implements Messages/DismissMessage
func (c *ZapClient) DismissMessage(id int) error {
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	resp, err := c.postJSON("Messages/DismissMessage", jsonStr, false)
	if err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		} else {
			return fmt.Errorf("dismissing the message wasn't successful")
		}
	}
	return err
}

//DismissAllMessages is equal to DismissMessage(-1)
func (c *ZapClient) DismissAllMessages() error {
	return c.DismissMessage(-1)
}

//GetAlertsTable implements Messages/GetAlertsTable with the default body
func (c *ZapClient) GetAlertsTable() (AlertsTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Title","name":"Title","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	resp, err := c.postJSON("Messages/GetAlertsTable", jsonStr, true)
	if err == nil {
		var alerts AlertsTable
		if json.Unmarshal(resp, &alerts) == nil {
			return alerts, nil
		} else {
			return *new(AlertsTable), err
		}
	}
	return *new(AlertsTable), err
}

//GetUnreadMessages partly implements Messages/Unread, it sets include_content & include_alerts to true
func (c *ZapClient) GetUnreadMessages() (UnreadMessages, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"Messages/Unread?include_content=true&include_alerts=true", nil)
	if err == nil {
		if token, err := c.GetNewToken(); err != nil {
			return *new(UnreadMessages), err
		} else {
			req.Header.Set("__RequestVerificationToken", token)
			if resp, err := c.client.Do(req); err == nil {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					var unread UnreadMessages
					if err := json.Unmarshal(body, &unread); err == nil {
						return unread, nil
					} else {
						return *new(UnreadMessages), err
					}
				}
			}
		}
	}
	return *new(UnreadMessages), err
}

//DismissAlert implements Messages/DismissAlert
func (c *ZapClient) DismissAlert(id int) error {
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	resp, err := c.postJSON("Messages/DismissAlert", jsonStr, false)
	if err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		} else {
			return fmt.Errorf("dismissing the alert failed")
		}
	}
	return err
}

//DismissAllAlerts is equal to DismissAlert(-1)
func (c *ZapClient) DismissAllAlerts() error {
	return c.DismissAlert(-1)
}

//SendMessage implements Messages/SendMessage. Doesn't return the ID yet, I didn't need it yet.
func (c *ZapClient) SendMessage(message string, toID uint) error {
	msg := chatMessage{Content: message, ID: toID, IsChat: true}
	jsonSlc, err := json.Marshal(msg)
	if err == nil {
		if resp, err := c.postJSON("Messages/SendMessage", string(jsonSlc), true); err == nil {
			if strings.HasPrefix(string(resp), `{"success":true,"result":"Success",`) {
				return nil
			} else {
				return fmt.Errorf("sending the message wasn't successful")
			}
		} else {
			return err
		}
	}
	return err
}
