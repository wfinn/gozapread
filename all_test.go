package gozapread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var zapread *ZapClient
var config TestConfig

type TestConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	if data, err := ioutil.ReadFile("./testconfig.json"); err == nil {
		if json.Unmarshal(data, &config) != nil {
			fmt.Println("There's an error in testconfig.json:", err)
			return
		}
	} else {
		fmt.Println("Couldn't read testconfig.json:", err)
	}
	var err error
	zapread, err = Login(config.Username, config.Password)
	if err != nil {
		fmt.Println("Login failed, no tests run.", err)
		return
	}
	os.Exit(m.Run())
}

func TestJoinGroup(t *testing.T) {
	if err := zapread.JoinGroup(199); err != nil {
		t.Error(err)
	}
}

func TestLeaveGroup(t *testing.T) {
	if err := zapread.LeaveGroup(199); err != nil {
		t.Error(err)
	}
}

func TestBalance(t *testing.T) {
	if balance, err := zapread.Balance(); err != nil {
		t.Error(err)
	} else {
		if balance < 1 {
			t.Error("The testing account should at least have 1 satoshi")
		}
	}
}

func TestGetUserID(t *testing.T) {
	userID, err := zapread.GetUserID("Zelgada")
	if err != nil {
		t.Error(err)
	}
	if userID != 1 {
		t.Error("Zelgada should be user 1")
	}
}

func TestIsOnline(t *testing.T) {
	var err error
	var idOnline bool
	if idOnline, err = zapread.IsUserIDOnline(1); err != nil {
		t.Error(err)
	}
	var nameOnline bool
	if nameOnline, err = zapread.IsUserNameOnline("Zelgada"); err != nil {
		t.Error(err)
	}
	if nameOnline != idOnline {
		t.Error("IsUserNameOnline and IsUserIdOnline had different results, might be a race condition")
	}
}

func TestGetNewToken(t *testing.T) {
	if token, err := zapread.GetNewToken(); err != nil {
		t.Error(err)
	} else {
		if len(token) != 151 {
			t.Error("GetNewToken didn't return a valid CSRF token")
		}
	}
}

func TestGetGroupID(t *testing.T) {
	if groupID := zapread.GetGroupID(157); groupID == 0 {
		t.Error("GetGroupId returned 0")
		//Maybe GetGroupId should return an error
	} else {
		if groupID != 28 {
			t.Error("GetGroupId returned a wrong ID. Expecated 28, got " + fmt.Sprint(groupID))
		}
	}
}

func TestGetMessageTable(t *testing.T) {
	if _, err := zapread.GetMessageTable(); err != nil {
		t.Error(err)
	}
}

func TestGetAlertsTable(t *testing.T) {
	if _, err := zapread.GetAlertsTable(); err != nil {
		t.Error(err)
	}
}

func TestGetUnreadMessages(t *testing.T) {
	if _, err := zapread.GetUnreadMessages(); err != nil {
		t.Error(err)
	}
}

func TestMessages(t *testing.T) {
	userID, _ := zapread.GetUserID(config.Username)
	if err := zapread.SendMessage("<b>test</b>", userID); err != nil {
		t.Error(err) //This is probably the only place where Error is correct, I usually want Log and FailNow
	}
	if err := zapread.DismissAllMessages(); err != nil {
		t.Error(err)
	}
}

func TestVoteComment(t *testing.T) {
	t.Skip("Voting costs money")
	if err := zapread.VoteComment(21764, true, 1); err != nil {
		t.Error(err)
	}
}

func TestVotePost(t *testing.T) {
	t.Skip("Voting costs money")
	if err := zapread.VotePost(1, true, 1); err != nil {
		t.Error(err)
	}
}

func TestUpdateAboutMe(t *testing.T) {
	newaboutme := fmt.Sprint(rand.Float64())
	if err := zapread.UpdateAboutMe(newaboutme); err != nil {
		t.Error(err)
	}
	if resp, err := zapread.client.Get(zapread.url + "user/" + config.Username); err != nil {
		t.Error(err)
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if !strings.Contains(string(body), newaboutme) {
				t.Error(fmt.Errorf("couldn't find the new about me text"))
			}
		} else {
			t.Error(err)
		}
	}
}
