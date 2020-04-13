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
	Balance string `json:"balance"`
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
		MessageID        uint   `json:"MessageId"`
		FromID           uint   `json:"FromId"`
		FromName         string `json:"FromName"`
		ToID             uint   `json:"ToId"`
		ToName           string `json:"ToName"`
		IsPrivateMessage bool   `json:"IsPrivateMessage"`
		TimeStamp        string `json:"TimeStamp"`
		Content          string `json:"Content"`
	} `json:"messages"`
}

type Invoice struct {
	Invoice string `json:"Invoice"`
	Result  string `json:"Result"`
	ID      uint   `json:"Id"`
}

type PaymentResp struct {
	Result string `json:"Result"`
	Fees   uint   `json:"Fees"`
}

type PaymentReq struct {
	Success     bool   `json:"success"`
	NumSatoshis string `json:"num_satoshis"`
	Destination string `json:"destination"`
}

type Tip struct {
	From    string
	Amount  uint
	AlertID uint
}

type ChatMessage struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	IsChat  bool   `json:"isChat"`
}

type UserHover struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

type PaymentCheck struct {
	Success bool `json:"success"`
	Result  bool `json:"result"`
}

type InvoiceResp struct {
	Invoice   string `json:"invoice"`
	IsDeposit bool   `json:"isDeposit"`
}
