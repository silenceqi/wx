package wx

import (
	"errors"
	"time"
)

type CacheItem struct {
	Value     string
	ExpiresIn int64
	UpdatedAt time.Time
}

func (item *CacheItem) IsExpired() bool {
	periodDuration := time.Duration((item.ExpiresIn - 120) * 1000 * 1000 * 1000)
	return item.UpdatedAt.Add(periodDuration).Before(time.Now())
}

type WXCache interface {
	GetAccessToken() (string, error)
	GetJsTicket() (string, error)
}

func NewWXCache(appId, secret string) WXCache {
	return DefaultWXCache{
		"appId":  appId,
		"secret": secret,
	}
}

type DefaultWXCache map[string]interface{}

func (cache DefaultWXCache) GetAccessToken() (string, error) {
	if accessTokenCache, ok := cache["accessTokenCache"]; ok {
		if cacheItem, ok := accessTokenCache.(*CacheItem); ok && !cacheItem.IsExpired() {
			return cacheItem.Value, nil
		}
	}
	weixin := NewWeixin(cache["appId"].(string), cache["secret"].(string))
	accessTokenRes, err := weixin.GetAccessToken()
	if err != nil {
		return "", err
	}
	if accessTokenRes.ErrCode != 0 {
		return "", errors.New(accessTokenRes.ErrMsg)
	}
	//cache accesstoken
	tokenCache := &CacheItem{
		Value:     accessTokenRes.AccessToken,
		UpdatedAt: time.Now(),
		ExpiresIn: accessTokenRes.ExpiresIn,
	}
	cache["accessTokenCache"] = tokenCache
	return accessTokenRes.AccessToken, nil
}

func (cache DefaultWXCache) GetJsTicket() (string, error) {
	if jsTicketCache, ok := cache["jsTicketCache"]; ok {
		if cacheItem, ok := jsTicketCache.(*CacheItem); ok && !cacheItem.IsExpired() {
			return cacheItem.Value, nil
		}
	}
	accessToken, err := cache.GetAccessToken()
	if err != nil {
		return "", err
	}
	wxJsSdk := NewJsSdk()
	ticketRes, err := wxJsSdk.GetTicket(accessToken)
	if err != nil {
		return "", err
	}
	if ticketRes.ErrCode != 0 {
		return "", errors.New(ticketRes.ErrMsg)
	}
	//cache jsTicket
	ticketCache := &CacheItem{
		Value:     ticketRes.Ticket,
		UpdatedAt: time.Now(),
		ExpiresIn: ticketRes.ExpiresIn,
	}
	cache["jsTicketCache"] = ticketCache
	return ticketCache.Value, nil

}
