package h2

import (
	"fmt"
	"net/url"
	"strings"

	tls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type Conn struct {
	Url            *url.URL
	Conn           *http2.Framer
	UnderlyingConn *tls.UConn
	Config         ReqConfig
	Client         *Client
	FirstUse       bool
}

// i suggest using webshare proxies if your using this library on a vps, since it will allow requests to go through
// if you dont use a proxy on a vps, you wont be able to send any requests due to some bot providers detecting
// datacenter ips.

// Connects to the url you supply, and stores it inside the Client struct.
func (Data *Client) Connect(addr string, config ReqConfig) (Connection Conn, err error) {
	Connection.Url = GrabUrl(addr)
	Connection.Client = Data
	Connection.Config = config

	if err := Connection.GenerateConn(config); err != nil {
		return Conn{}, err
	}

	return Connection, nil
}

// Does a request, since http2 doesnt like to resent new headers. after the first request it will reconnect to the server
// and make a new http2 framer variable to use.
func (Data *Conn) Do(method, json, content_type string, cookies *[]string) (Config Response, err error) {
	if !Data.FirstUse {
		Data.FirstUse = true
	} else {
		if err = Data.GenerateConn(Data.Config); err != nil {
			return Response{}, err
		}
	}

	if cookies != nil {
		Data.Client.Config.Headers["cookie"] += TurnCookieHeader(*cookies)
	}

	var FoundAndSent bool

	Headers := Data.GetHeaders(method)

	if method != "GET" {
		var FoundType, FoundLength, FoundEither bool
		for _, header := range Headers {
			if strings.Contains(header, "content-type") {
				FoundType = true
			}
			if strings.Contains(header, "content-length") {
				FoundLength = true
			}
		}
		if !FoundLength {
			Data.AddHeader("content-length", fmt.Sprintf("%v", len(json)))
			FoundEither = true
		}
		if !FoundType {
			Data.AddHeader("content-type", content_type)
			FoundEither = true
		}
		if FoundEither {
			Headers = Data.GetHeaders(method)
			Data.SendHeaders(Headers, method == "GET")
			FoundAndSent = true
		}
		Data.DataSend([]byte(json))
	}
	if !FoundAndSent {
		Data.SendHeaders(Headers, method == "GET")
	}
	resp, err := Data.FindData()
	if err != nil {
		return resp, err
	}
	if resp.Status == "302" || resp.Status == "301" {
		Data.Url.Path = GetHeaderVal("location", resp.Headers).Value
		if err = Data.GenerateConn(Data.Config); err != nil {
			return Response{}, err
		}
		Data.SendHeaders(Data.GetHeaders(method), method == "GET")
		return Data.FindData()
	}
	return resp, err
}

func (Data *Conn) ChangeProxy(proxy *ProxyAuth) {
	Data.Config.Proxy = proxy
}

// Changes the url path, so you can send to different locations under one variable.
func (Data *Conn) ChangeURL(url *url.URL) {
	Data.Url = url
}

// adds a header to the client struct
func (Data *Conn) AddHeader(name, value string) {
	Data.Client.Config.Headers[name] = value
}

// deletes headers from a client struct
func (Data *Conn) DeleteHeader(headernames ...string) {
	for _, val := range headernames {
		delete(Data.Client.Config.Headers, val)
	}
}
