package gowechat

import (
	"fmt"
)

type Wechat struct {
	Token          string
	EncodingAESKey string
	CorpId         string
	CorpSecret     string
}

func New(token, encodingAesKey, corpId, corpSecret string) *Wechat {
	return &Wechat{
		Token:          token,
		EncodingAESKey: encodingAesKey,
		CorpId:         corpId,
		CorpSecret:     corpSecret,
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
