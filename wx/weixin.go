package wx

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json: "expires_in"`
}

type Weixin interface {
	GetAccessToken() (*AccessTokenResponse, error)
}

var appInfo struct {
	appid, secret string
	isInit        bool
}

func Init(appid, secret string) Weixin {
	appInfo.appid = appid
	appInfo.secret = secret
	appInfo.isInit = true
	return NewWeixin(appid, secret)
}

func NewWeixin(appid, secret string) Weixin {
	if appid == "" && appInfo.isInit {
		appid = appInfo.appid
	}
	if secret == "" && appInfo.isInit {
		secret = appInfo.secret
	}
	return &defaultWeixin{
		AppId:  appid,
		Secret: secret,
	}
}

type defaultWeixin struct {
	AppId  string
	Secret string
}

func (weixin *defaultWeixin) GetAccessToken() (*AccessTokenResponse, error) {
	var serviceUrl string = "https://api.weixin.qq.com/cgi-bin/token"
	reqParams := url.Values{}
	reqParams.Add("grant_type", "client_credential")
	reqParams.Add("appid", weixin.AppId)
	reqParams.Add("secret", weixin.Secret)
	serviceUrl += "?" + reqParams.Encode()
	res, err := http.Get(serviceUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	tokenRes := &AccessTokenResponse{}
	err = json.Unmarshal(resBuf, tokenRes)
	if err != nil {
		return nil, err
	}
	return tokenRes, nil
}
