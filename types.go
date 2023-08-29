package gowechat

type WechatUserMsg struct {
	ToUserName string `xml:"ToUserName"`
	CreateTime int64  `xml:"CreateTime"`
	MsgType    string `xml:"MsgType"`
	Event      string `xml:"Event"`
	Token      string `xml:"Token"`
	OpenKfId   string `xml:"OpenKfId"`
}

type AccessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type MsgRet struct {
	Errcode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	NextCursor string `json:"next_cursor"`
	MsgList    []Msg  `json:"msg_list"`
}

type MsgEvent struct {
	EventType      string `json:"event_type"`
	Scene          string `json:"scene"`
	OpenKfid       string `json:"open_kfid"`
	ExternalUserid string `json:"external_userid"`
	WelcomeCode    string `json:"welcome_code"`
}

type Msg struct {
	MsgId    string   `json:"msgid"`
	SendTime int64    `json:"send_time"`
	Origin   int      `json:"origin"`
	Msgtype  string   `json:"msgtype"`
	Event    MsgEvent `json:"event"`
	Text     struct {
		Content string `json:"content"`
	} `json:"text"`
	Link struct {
		Title  string `json:"title"`
		Desc   string `json:"desc"`
		Url    string `json:"url"`
		PicUrl string `json:"pic_url"`
	} `json:"link"`
	OpenKfid       string `json:"open_kfid"`
	ExternalUserid string `json:"external_userid"`
}

type ReplyMsg struct {
	Touser   string `json:"touser,omitempty"`
	OpenKfid string `json:"open_kfid,omitempty"`
	MsgId    string `json:"msgid,omitempty"`
	MsgType  string `json:"msgtype,omitempty"`
	Text     struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
}
