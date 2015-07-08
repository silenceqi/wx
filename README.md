# Go-Webchat-JSSDK lib
simple library for weixin jssdk signature

# Useage

package main

import (
	"log"
	"text/template"

	"net/http"

	"wx"
)

func checkRuntimeErr(err error) bool {
	if err != nil {
		log.Fatalf("error: %v", err)
		return true
	}
	return false
}

var globalCache wx.DefaultWXCache = wx.DefaultWXCache{
	"appId":  "YOUR APPID",
	"secret": "YOUR APP SECRET",
}

func main() {
	accessToken, err := globalCache.GetAccessToken()
	checkRuntimeErr(err)
	ticket, err := globalCache.GetJsTicket()
	checkRuntimeErr(err)
	wxJsSdk := wx.NewJsSdk()
	http.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		var signURL = req.URL.Host + req.URL.RequestURI()
		spkg, _ := wxJsSdk.ComputeSignature(accessToken, ticket, signURL)
		spkg["appId"] = globalCache["appId"]
		spkg["baseuri"] = req.URL.RequestURI()
		log.Println(spkg)
		tpl, err := template.ParseFiles("./templates/index.html")
		checkRuntimeErr(err)
		tpl.Execute(w, spkg)
	})
	err = http.ListenAndServe(":9090", nil)
	checkRuntimeErr(err)
}
