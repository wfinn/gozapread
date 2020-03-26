# gozapread
[https://zapread.com](https://zapread.com)

## Why?
I had an idea for a zapread bot and therefore needed this.

## How it works

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

### Tips
- You can use *go doc gozapread ZapClient* to list all available functionality.
- GetUnreadMessages() + DismissMessage()

# Todo
- refactor the hell out of this mess
- provide more examples

## Next endpoints, probably
- **ChangePassword(old, new string) error** /Manage/ChangePassword
- /Vote/Comment like /Vote/Post
- /Home/TopPosts
- /Manage/UpdateAboutMe/ __RequestVerificationToken=abc&AboutMe=text
