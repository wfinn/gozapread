# gozapread
A code base for [zapread](https://github.com/Horndev/zapread.com) bots.  


## Why?
I had an idea for a zapread bot and therefore needed this.  
Also, I am not the first person that needed something like this, I hope this is useful for someone.  

## How it works
The `Login` function will create a `ZapClient` for you. This `ZapClient` will manage the session for you. Sessions don't expire quickly, just make sure to fire a request every now and then so your session cookie can be refreshed.


Then you can call functions, the names are similar to the actual endpoints (e.g. /Account/Balance is `Balance()`).  
As far as I can tell, most important functionality is already implemented, see `go doc` for more details. If something is missing, feel free to create an issue or PR.

Most endpoints accept and return simple structured json, so there are many custom structs. Some things are not easily possible by calling zapreads API directly, so there are some convinience functions like `ParseTips`.
### Example
The following few lines would make a new post.
```go
api, err := gozapread.Login("bot", "password")
if err != nil {
	log.Fatal("Login failed.")
}
api.SubmitNewPost("New Post", "Hi, I am a <b>bot</b>!", 199)
```
Let's do something with the data `SubmitNewPost` returned.
```go
if resp, err := api.SubmitNewPost("New Post", "Hi, I am a <b>bot</b>!", 199); err == nil {
	fmt.Printf("New post with id %d", resp.PostID)
}
```
With this simple setup you could create bots that report some statistics about Bitcoin for example. I'm sure that would earn you some satoshi.

## Tests
Before running `go test` you need to put the login details of a zapread account in testconfig.json.
(Then run `git update-index --skip-worktree testconfig.json` to ignore the change.)

Currently 1 satoshi is required in the test account, more requirements may follow with other tests.

## Todo
- provide more examples

## Next endpoints, probably
- **ChangePassword(old, new string) error** /Manage/ChangePassword
- /Vote/Comment like /Vote/Post
- /Home/TopPosts
- /Manage/UpdateAboutMe/ __RequestVerificationToken=abc&AboutMe=text
