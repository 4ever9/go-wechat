package gowechat

type WechatUserAskMsg struct {
	ToUserName string `xml:"ToUserName"`
	CreateTime int64  `xml:"CreateTime"`
	MsgType    string `xml:"MsgType"`
	Event      string `xml:"Event"`
	Token      string `xml:"Token"`
	OpenKfId   string `xml:"OpenKfId"`
}

type AccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type MsgRet struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	NextCursor string `json:"next_cursor"`
	MsgList    []Msg  `json:"msg_list"`
}

type Msg struct {
	Msgid    string `json:"msgid"`
	SendTime int64  `json:"send_time"`
	Origin   int    `json:"origin"`
	Msgtype  string `json:"msgtype"`
	Event    struct {
		EventType      string `json:"event_type"`
		Scene          string `json:"scene"`
		OpenKfid       string `json:"open_kfid"`
		ExternalUserid string `json:"external_userid"`
		WelcomeCode    string `json:"welcome_code"`
	} `json:"event"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	OpenKfid       string `json:"open_kfid"`
	ExternalUserid string `json:"external_userid"`
}

type ReplyMsg struct {
	Touser   string `json:"touser,omitempty"`
	OpenKfid string `json:"open_kfid,omitempty"`
	Msgid    string `json:"msgid,omitempty"`
	Msgtype  string `json:"msgtype,omitempty"`
	Text     struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
}
