package gozapread

type MessageTable struct {
	Draw            uint `json:"draw"`
	RecordsTotal    uint `json:"recordsTotal"`
	RecordsFiltered uint `json:"recordsFiltered"`
	Data            []struct {
		ID      uint   `json:"Id"`
		Status  string `json:"Status"`
		Type    string `json:"Type"`
		From    string `json:"From"`
		FromID  string `json:"FromID"`
		Date    string `json:"Date"`
		Link    string `json:"Link"`
		Anchor  string `json:"Anchor"`
		Message string `json:"Message"`
	} `json:"data"`
}

type Post struct {
	PostID   uint   `json:"PostId"`
	Content  string `json:"Content"`
	GroupID  uint   `json:"GroupId"`
	UserID   bool   `json:"UserId"`
	Title    string `json:"Title"`
	IsDraft  bool   `json:"IsDraft"`
	Language string `json:"Language"`
}

type PostResp struct {
	Result      string `json:"result"`
	Success     bool   `json:"success"`
	PostID      uint   `json:"postId"`
	HTMLContent string `json:"HTMLContent"`
}

type Comment struct {
	CommentContent string `json:"CommentContent"`
	PostID         uint   `json:"PostId"`
	CommentID      uint   `json:"CommentId"`
	IsReply        bool   `json:"IsReply"`
}

type BalanceResp struct {
	Balance uint `json:"balance"`
}

type AlertsTable struct {
	Draw            uint `json:"draw"`
	RecordsTotal    uint `json:"recordsTotal"`
	RecordsFiltered uint `json:"recordsFiltered"`
	Data            []struct {
		AlertID        uint   `json:"AlertId"`
		Status         string `json:"Status"`
		Title          string `json:"Title"`
		Date           string `json:"Date"`
		Link           string `json:"Link"`
		Anchor         string `json:"Anchor"`
		Message        string `json:"Message"`
		HasCommentLink bool   `json:"HasCommentLink"`
		HasLink        bool   `json:"HasLink"`
	} `json:"data"`
}

type UnreadMessages struct {
	Success  bool `json:"success"`
	Messages []struct {
		MessageID        int    `json:"MessageId"`
		FromID           int    `json:"FromId"`
		FromName         string `json:"FromName"`
		ToID             int    `json:"ToId"`
		ToName           string `json:"ToName"`
		IsPrivateMessage bool   `json:"IsPrivateMessage"`
		TimeStamp        string `json:"TimeStamp"`
		Content          string `json:"Content"`
	} `json:"messages"`
}
