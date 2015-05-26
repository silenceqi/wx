package wx

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type TicketResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

type SignPackage map[string]interface{}
type JsSdk interface {
	GenerateNoncestr(length int) string
	GetTicket(token string) (*TicketResponse, error)
	ComputeSignature(token, ticket, signUrl string) (SignPackage, error)
}

type JsSdkDefault struct {
}

func NewJsSdk() JsSdk {
	return &JsSdkDefault{}
}

func (sdk *JsSdkDefault) GenerateNoncestr(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var str string
	var ridx int
	upBound := len(chars)
	for i := 0; i < length; i++ {
		ridx = rand.Intn(upBound)
		str += chars[ridx : ridx+1]
	}
	return str
}

func (sdk *JsSdkDefault) GetTicket(accessToken string) (*TicketResponse, error) {
	serviceUrl := "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
	reqParams := url.Values{}
	reqParams.Add("type", "jsapi")
	reqParams.Add("access_token", accessToken)
	serviceUrl += "?" + reqParams.Encode()
	res, err := http.Get(serviceUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytesBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ticketRes := &TicketResponse{}
	err = json.Unmarshal(bytesBuf, ticketRes)
	if err != nil {
		return nil, err
	}
	return ticketRes, nil
}

func (sdk *JsSdkDefault) ComputeSignature(accessToken, ticket, signUrl string) (SignPackage, error) {
	timestamp := time.Now().Unix()
	noncestr := sdk.GenerateNoncestr(16)
	params := PkgComponents{
		"jsapi_ticket": ticket,
		"noncestr":     noncestr,
		"timestamp":    fmt.Sprintf("%d", timestamp),
		"url":          signUrl,
	}
	rawStr := params.buildRawQueryString()
	signature := sha1.Sum([]byte(rawStr))

	return SignPackage{
		"access_token": accessToken,
		"ticket":       ticket,
		"noncestr":     noncestr,
		"timestamp":    timestamp,
		"signature":    fmt.Sprintf("%x", signature),
	}, nil
}

type PkgComponents map[string]string

func (pc PkgComponents) buildRawQueryString() string {
	qstr := ""
	keys := []string{}
	for key := range pc {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		qstr += key + "=" + pc[key] + "&"
	}
	return qstr[0 : len(qstr)-1]
}
