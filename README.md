# gozapread
Some go code I need for building bots for [https://zapread.com](https://zapread.com).

## Why?
I had an idea for a zapread bot and therefore needed this.

## How it works (How it will work)

### Example
The ```gozapread.Login(user, pass string)``` function prepares an internal http.Client which handles the session for you.

```go
api, err := gozapread.Login("bot", "password")
if err != nil {
	log.Fatal("Login failed.")
}
api.SubmitNewPost("New Post", "Hi, I am a <b>bot</b>!", 199)
```
zapread.com usually returns some json, I use [https://mholt.github.io/json-to-go/](https://mholt.github.io/json-to-go/) to prepare structs this data can be parsed to.
Let's do something with the data SubmitNewPost returned.
```go
if resp, err := api.SubmitNewPost("New Post", "Hi, I am a <b>bot</b>!", 199); err == nil {
	fmt.Printf("New post with id %d", resp.PostID)
}
```
### Functions
- **Login(user, pass string) (zapclient, error)** prepares client
- **AddComment(content string, postid, commentid uint) error**
- **DismissAlert(id uint) error** // should be int -1 means dismiss all
- **DismissMessage(id uint) error** // should be int -1 means dismiss all
- **GetAlertsTable() (AlertsTable, error)**
- **GetGroupId(postid uint) (result uint)** // return an error
- **GetMessageTable() (MessageTable, error)**
- **GetNewToken() (string, error)**
- **JoinGroup(groupid uint) error**
- **LeaveGroup(groupid uint) error**
- **SubmitNewPost(title, content string, groupid uint) (PostResp, error)**
- **TipUser(userid, amount uint) error**
- **UnreadMessages() bool** //TODO return the uint instead
- **UserBalance() (uint, error)**
- **VotePost(postid int, upvote bool, amount uint) error**

#### Not Implemented
- **ChangePassword(old, new string) error** /Manage/ChangePassword
- /Messages/SendMessage/
- /Vote/Comment like /Vote/Post
- /Home/TopPosts
- /Messages/DismissAlert -1 means all
- /Lightning/GetDepositInvoice/ {"amount":"1","memo":"ZapRead.com deposit","anon":"0","use":"userDeposit","useId":-1,"useAction":-1} {"Invoice":"blah","Result":"success","Id":123456}
- /Lightning/ValidatePaymentRequest {"request":"blah"} {"success":true,"num_satoshis":"1","destination":"blah"}
- /Lightning/SubmitPaymentRequest {"request":"blah"}
- /Manage/UpdateAboutMe/ __RequestVerificationToken=abc&AboutMe=text

# Todo
- generate function list
- refactor the hell out of this mess
