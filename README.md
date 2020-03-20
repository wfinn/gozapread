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
- **Login(user, pass string) error** prepares client
- **extractToken(html string) string** can extract a token from a response body
- TODO **GetNewToken(path string) (string, err)** Refactoring (See if path is nessecary, probably not)
- **UnreadMessages() uint** /Messages/UnreadMessages/
- **GetMessageTable()** /Messages/GetMessagesTable Unmarshal the json to a MessageTable
- **DismissMessage(id uint)** /Messages/DismissMessage
- **GetGroupId(postId uint) uint** curl --silent https://www.zapread.com/Post/Detail/6126 | grep -o "data-groupid=\"[^\"]\*" | sed 's/^data-groupid="//'
- **SubmitNewPost(title, body string, group uint) bool** /Post/SubmitNewPost/
- TODO **ChangePassword(old, new string) error** /Manage/ChangePassword

# Todo
- list more functions that need to be implemented (I should probably parse them from zapread.com code...)
- refactor the hell out of this mess
