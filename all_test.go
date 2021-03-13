package gozapread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var zapread *ZapClient

type TestConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestMain(m *testing.M) {
	var config TestConfig
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

func TestGetUserId(t *testing.T) {
	userID, err := zapread.GetUserId("Zelgada")
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
	if idOnline, err = zapread.IsUserIdOnline(1); err != nil {
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

func TestGetGroupId(t *testing.T) {
	if groupID := zapread.GetGroupId(157); groupID == 0 {
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
