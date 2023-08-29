package gowechat

import (
	"encoding/xml"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	tokenExpiration = 5 * time.Minute
	retryExpiration = 1 * time.Minute
	retryCacheKey   = "wecom:retry:%s"
	tokenCacheKey   = "wecom:token"
)

var tokenCache = cache.New(tokenExpiration, tokenExpiration)
var retryCache = cache.New(retryExpiration, retryExpiration)

type Wechat struct {
	Token          string
	EncodingAESKey string
	CorpId         string
	CorpSecret     string
	httpClient     *req.Client
	nextCursor     string
}

func New(token, encodingAesKey, corpId, corpSecret string) *Wechat {
	return &Wechat{
		Token:          token,
		EncodingAESKey: encodingAesKey,
		CorpId:         corpId,
		CorpSecret:     corpSecret,
		httpClient:     req.C(),
	}
}

func (we *Wechat) CheckWeComSign(msgSign, timestamp, nonce, echoStr string) ([]byte, error) {
	bizCrypto := NewWXBizMsgCrypt(we.Token, we.EncodingAESKey, we.CorpId, 1)
	decoded, err := bizCrypto.VerifyURL(msgSign, timestamp, nonce, echoStr)
	if nil != err {
		return nil, fmt.Errorf("code: %d, msg: %s", err.ErrCode, err.ErrMsg)
	}

	return decoded, nil
}

func isRetry(signature string) bool {
	key := fmt.Sprintf(retryCacheKey, signature)
	_, found := retryCache.Get(key)
	if found {
		return true
	}

	retryCache.Set(key, "true", retryExpiration)
	return false
}

func (we *Wechat) DecryptWeComMsg(msgSign, timestamp, nonce string, body []byte) (*MsgRet, error) {
	crypt := NewWXBizMsgCrypt(we.Token, we.EncodingAESKey, we.CorpId, 1)
	data, cryptoErr := crypt.DecryptMsg(msgSign, timestamp, nonce, body)
	if cryptoErr != nil {
		return nil, fmt.Errorf("decrypt message: %w", cryptoErr)
	}
	var wechatUserMsg WechatUserMsg
	err := xml.Unmarshal(data, &wechatUserMsg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	accessToken, err := we.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	msgToken := wechatUserMsg.Token
	msgRet, err := we.getMsgs(accessToken, msgToken)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	we.nextCursor = msgRet.NextCursor

	if isRetry(msgSign) {
		return nil, nil
	}

	return msgRet, nil
}

func (we *Wechat) getAccessToken() (string, error) {
	data, found := tokenCache.Get(tokenCacheKey)
	if found {
		return fmt.Sprintf("%v", data), nil
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		we.CorpId,
		we.CorpSecret,
	)

	var accessToken AccessToken
	_, err := we.httpClient.R().SetSuccessResult(&accessToken).Get(url)
	if err != nil {
		return "", fmt.Errorf("get token by url: %w", err)
	}

	token := accessToken.AccessToken
	tokenCache.Set(tokenCacheKey, token, tokenExpiration)

	return token, nil
}

func (we *Wechat) getMsgs(accessToken, msgToken string) (*MsgRet, error) {
	var msgRet *MsgRet
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token=" + accessToken
	args := map[string]string{
		"token": msgToken,
	}
	if we.nextCursor != "" {
		args["cursor"] = we.nextCursor
	}
	_, err := we.httpClient.R().
		SetBody(args).
		SetSuccessResult(&msgRet).
		Post(url)
	if err != nil {
		return nil, err
	}

	return msgRet, nil
}

func (we *Wechat) ReplyMessage(externalUserId, openKfid, ask, content string) error {
	reply := ReplyMsg{
		Touser:   externalUserId,
		OpenKfid: openKfid,
		MsgType:  "text",
		Text: struct {
			Content string `json:"content,omitempty"`
		}{Content: content},
	}

	token, err := we.getAccessToken()
	if err != nil {
		return err
	}

	return we.callTalk(reply, token)
}

func (we *Wechat) callTalk(reply ReplyMsg, accessToken string) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/kf/send_msg?access_token=" + accessToken
	_, err := we.httpClient.R().SetBody(reply).Post(url)

	if err != nil {
		return err
	}

	return nil
}
