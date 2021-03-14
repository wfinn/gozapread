package gozapread

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c *ZapClient) UnreadMessages() bool { //TODO return the uint instead
	if res, err := c.client.Get(c.url + "Messages/UnreadMessages/"); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return !(string(body) == "0")
		}
	}
	return false
}

func (c *ZapClient) GetMessageTable() (MessageTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"From","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Message","name":"Message","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJSON("Messages/GetMessagesTable", jsonStr, true); err == nil {
		var messages MessageTable
		if json.Unmarshal(resp, &messages) == nil {
			return messages, nil
		}
	}
	return *new(MessageTable), errors.New("GetMessageTable failed")
}

func (c *ZapClient) DismissMessage(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJSON("Messages/DismissMessage", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}
	}
	return errors.New("DismissMessage failed")
}

func (c *ZapClient) GetAlertsTable() (AlertsTable, error) {
	jsonStr := `{"draw":1,"columns":[{"data":null,"name":"Status","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Date","name":"Date","searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"Title","name":"Title","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Link","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":null,"name":"Action","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":1,"dir":"desc"}],"start":0,"length":25,"search":{"value":"","regex":false}}`
	if resp, err := c.postJSON("Messages/GetAlertsTable", jsonStr, true); err == nil {
		var alerts AlertsTable
		if json.Unmarshal(resp, &alerts) == nil {
			return alerts, nil
		}
	}
	return *new(AlertsTable), errors.New("GetAlertsTable failed")
}

func (c *ZapClient) GetUnreadMessages() (UnreadMessages, error) {
	if req, err := http.NewRequest(http.MethodGet, c.url+"Messages/Unread?include_content=true&include_alerts=true", nil); err == nil {
		if token, err := c.GetNewToken(); err == nil {
			req.Header.Set("__RequestVerificationToken", token)
			if resp, err := c.client.Do(req); err == nil {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					var unread UnreadMessages
					if json.Unmarshal(body, &unread) == nil {
						return unread, nil
					}
				}
			}
		}
	}
	return *new(UnreadMessages), errors.New("UnreadMessages failed")
}

func (c *ZapClient) DismissAlert(id uint) error { // should be int -1 means dismiss all
	jsonStr := fmt.Sprintf(`{"id":%d}`, id)
	if resp, err := c.postJSON("Messages/DismissAlert", jsonStr, false); err == nil {
		if string(resp) == `{"Result":"Success"}` {
			return nil
		}

	}
	return errors.New("DismissAlert failed")
}

//Sends a private message, doesn't return the ID yet, I didn't need it yet.
func (c *ZapClient) SendMessage(message string, toID uint) error {
	msg := ChatMessage{Content: message, ID: toID, IsChat: true}
	if jsonSlc, err := json.Marshal(msg); err == nil {
		if resp, err := c.postJSON("Messages/SendMessage", string(jsonSlc), true); err == nil {
			if strings.HasPrefix(string(resp), `{"success":true,"result":"Success",`) {
				return nil
			}
		}
	}
	return errors.New("AddComment failed")
}