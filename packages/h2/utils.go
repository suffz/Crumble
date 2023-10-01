package h2

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	tls "github.com/bogdanfinn/utls"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func (Data *Conn) Connect(config ReqConfig) (net.Conn, bool, error) {
	if conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", config.Proxy.IP, config.Proxy.Port)); err == nil {
		Conn_URL := CheckAddr(Data.Url)
		if config.Proxy.User != "" && config.Proxy.Password != "" {
			conn.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Authorization: Basic %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\n\r\n", Conn_URL, Conn_URL, base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config.Proxy.User, config.Proxy.Password))))))
		} else {
			conn.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\n\r\n", Conn_URL, Conn_URL)))
		}
		var junk = make([]byte, 4096)
		conn.Read(junk)
		switch Status := string(junk); Status[9:12] {
		case "200":
			return conn, true, nil
		case "407":
			return nil, false, errors.New(fmt.Sprintf("[%v] Proxy <%v> Failed to authorize: Username/Password invalid.\n", Status[9:12], config.Proxy.IP))
		default:
			return nil, false, errors.New(fmt.Sprintf("[%v] Proxy <%v> Failed to authorize: Unknown Statuscode.\n", Status[9:12], config.Proxy.IP))
		}
	}
	return nil, false, errors.New(fmt.Sprintf("Proxy <%v> Failed to authorize: Unknown Error EOF\n", config.Proxy.IP))
}

// Generate conn performs a conn to the url you supply.
func (Data *Conn) GenerateConn(config ReqConfig) (err error) {
	if config.DataBodyMaxLength == 0 {
		return errors.New("error: [DataBodyMaxLength] cannot be 0, suggested value to be above 130000")
	}

	var conn net.Conn
	var tlsConn *tls.UConn
	if config.Proxy != nil {
		var ok bool
		var err error
		conn, ok, err = Data.Connect(config)
		if !ok && err != nil {
			return err
		}
	} else {
		conn, err = net.Dial("tcp", CheckAddr(Data.Url))
		if err != nil {
			return err
		}
	}

	tlsConn = tls.UClient(conn, &tls.Config{
		ServerName:               Data.Url.Host,
		NextProtos:               Data.Client.Config.Protocols,
		InsecureSkipVerify:       config.InsecureSkipVerify,
		Renegotiation:            config.Renegotiation,
		PreferServerCipherSuites: config.PreferServerCipherSuites,
		RootCAs:                  config.RootCAs,
		ClientCAs:                config.ClientCAs,
	}, config.BuildID, true, config.UseHTTP1)

	if config.SaveCookies {
		if Data.Client.Cookies == nil || len(Data.Client.Cookies) == 0 {
			Data.Client.Cookies = make(map[string][]hpack.HeaderField)
		}
	}

	fmt.Fprintf(tlsConn, http2.ClientPreface)

	if err = tlsConn.Handshake(); err != nil {
		return err
	}

	Data.Conn = http2.NewFramer(tlsConn, tlsConn)
	Data.Conn.SetReuseFrames()
	Data.WriteSettings(config.DataBodyMaxLength)
	Data.Windows_Update()
	Data.Send_Prio_Frames()
	Data.UnderlyingConn = tlsConn
	return nil
}

// gets a selected cookie based on the cookie_name variable
//
//	e.g. "__vf_bm" > "__vf_bm=awdawd223reqfqh32rqrf32qr"
func (Data *Conn) GetCookie(cookie_name, url string) string {
	for _, val := range Data.Client.Cookies[url] {
		if strings.Contains(val.Value, cookie_name) {
			Cookie := strings.Split(val.Value, "=")
			return fmt.Sprintf("%v=%v", Cookie[0], Cookie[1])
		}
	}

	return ""
}

// Gets a header value based on the name you supply.
func GetHeaderVal(name string, headers []hpack.HeaderField) hpack.HeaderField {
	for _, data := range headers {
		if data.Name == name {
			return data
		}
	}
	return hpack.HeaderField{}
}

// This is a helper function that gets all the cookies from a
// cached url and returns them in a format that works with the cookie: header.
func (Data *Conn) TransformCookies(url string) string {
	var cookies []string
	for _, val := range Data.Client.Cookies[url] {
		cookie_name := strings.Split(strings.Split(val.Value, ";")[0], "=")
		cookies = append(cookies, fmt.Sprintf("%v=%v", cookie_name[0], cookie_name[1]))
	}
	return strings.Join(cookies, "; ")
}

// strings.Join shortcut to turn your list of coookies into a cookie: header format.
func TurnCookieHeader(Cookies []string) string {
	return strings.Join(Cookies, "; ")
}

// Sends data through the framer
func (Data *Conn) DataSend(body []byte) {
	Data.Conn.WriteData(uint32(Data.Config.ID), true, body)
}

// Sends priority frames, this ensures the right data is sent in the correct order.
func (Data *Conn) Send_Prio_Frames() {
	Data.Conn.WritePriority(3, http2.PriorityParam{
		StreamDep: 0,
		Weight:    200,
		Exclusive: false,
	})

	Data.Conn.WritePriority(5, http2.PriorityParam{
		StreamDep: 0,
		Weight:    100,
		Exclusive: false,
	})

	Data.Conn.WritePriority(7, http2.PriorityParam{
		StreamDep: 0,
		Weight:    0,
		Exclusive: false,
	})

	Data.Conn.WritePriority(9, http2.PriorityParam{
		StreamDep: 7,
		Weight:    0,
		Exclusive: false,
	})

	Data.Conn.WritePriority(11, http2.PriorityParam{
		StreamDep: 3,
		Weight:    0,
		Exclusive: false,
	})

	Data.Conn.WritePriority(13, http2.PriorityParam{
		StreamDep: 0,
		Weight:    240,
		Exclusive: false,
	})
}

// Loops over the Config headers and applies them to the Client []string variable.
// Method for example "GET".
func (Data *Conn) GetHeaders(method string) (headers []string) {
	for _, name := range Data.Client.Config.HeaderOrder {
		switch name {
		case ":authority":
			headers = append(headers, name+": "+Data.Url.Host)
		case ":method":
			headers = append(headers, name+": "+method)
		case ":path":
			headers = append(headers, name+": "+CheckQuery(Data.Url))
		case ":scheme":
			headers = append(headers, name+": "+Data.Url.Scheme)
		default:
			if val, exists := Data.Client.Config.Headers[name]; exists {
				headers = append(headers, name+": "+val)
			}
		}
	}

	for name, val := range Data.Client.Config.Headers {
		if !strings.Contains(strings.Join(Data.Client.Config.HeaderOrder, ","), name) {
			headers = append(headers, name+": "+val)
		}
	}

	return
}

// Writes the headers to the http2 framer.
// this function also encodes the headers into a []byte
// Endstream is also called in this function, only use true values when performing GET requests.
func (Data *Conn) SendHeaders(headers []string, endStream bool) {
	Data.Conn.WriteHeaders(
		http2.HeadersFrameParam{
			StreamID:      uint32(Data.Config.ID),
			BlockFragment: Data.FormHeaderBytes(headers),
			EndHeaders:    true,
			EndStream:     endStream,
		},
	)
}

// Writes the window update frame to the http2 framer.
func (Data *Conn) Windows_Update() {
	Data.Conn.WriteWindowUpdate(0, 12517377)
}

// Write settings writes the default chrome settings to the framer
func (Data *Conn) WriteSettings(ResponseDataSize uint32) {
	Data.Conn.WriteSettings(
		http2.Setting{
			ID: http2.SettingHeaderTableSize, Val: 65536,
		},
		http2.Setting{
			ID: http2.SettingEnablePush, Val: 1,
		},
		http2.Setting{
			ID: http2.SettingMaxConcurrentStreams, Val: 1000,
		},
		http2.Setting{
			ID: http2.SettingInitialWindowSize, Val: ResponseDataSize,
		},
		http2.Setting{
			ID: http2.SettingMaxFrameSize, Val: 16384,
		},
		http2.Setting{
			ID: http2.SettingMaxHeaderListSize, Val: 262144,
		},
	)
}

// Find data is called after the prior settings/window/prio frames are performed, it goes through the
// framer and returns its data, any errors and also headers / status codes.
func (Datas *Conn) FindData() (Config Response, err error) {
	for {
		f, err := Datas.Conn.ReadFrame()
		if err != nil {
			return Config, err
		}
		switch f := f.(type) {
		case *http2.DataFrame:
			Config.Data = append(Config.Data, f.Data()...)
			if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) {
				return Config, nil
			}
		case *http2.HeadersFrame:
			Config.Headers, err = hpack.NewDecoder(Datas.Config.DataBodyMaxLength, nil).DecodeFull(f.HeaderBlockFragment())
			if err != nil {
				return Config, err
			}
			for _, Data := range Config.Headers {
				switch Data.Name {
				case ":status":
					Config.Status = Data.Value
				case "set-cookie":
					if Datas.Config.SaveCookies {
						Datas.Client.Cookies[Datas.Url.String()] = append(Datas.Client.Cookies[Datas.Url.String()], Data)
					}
				}
			}
			if f.FrameHeader.Flags.Has(http2.FlagDataEndStream) && f.FrameHeader.Flags.Has(http2.FlagHeadersEndStream) {
				return Config, nil
			}
		case *http2.RSTStreamFrame:
			return Config, errors.New(f.ErrCode.String())
		case *http2.GoAwayFrame:
			return Config, errors.New(f.ErrCode.String())
		}
	}
}

// Turns the addr into a url.URL variable.
func GrabUrl(addr string) *url.URL {
	URL, _ := url.Parse(addr)
	if URL.Path == "" {
		URL.Path = "/"
	}
	return URL
}

// Checks if there are params in your url and adds it to your path.
//
//	e.g. "/api/name?code=12343&scope=1234"
func CheckQuery(Data *url.URL) string {
	if Data.Query().Encode() != "" {
		return Data.Path + "?" + Data.Query().Encode()
	}
	return Data.Path
}

// Form header bytes takes the []string of headers and turns it into []byte data
// this is so it can be compatiable for the http2 headers.
func (Data *Conn) FormHeaderBytes(headers []string) []byte {
	var val []string
	hbuf := bytes.NewBuffer([]byte{})
	encoder := hpack.NewEncoder(hbuf)
	for _, header := range headers {
		switch data := strings.Split(header, ":"); len(data) {
		case 3:
			val = data[1:]
			val[0] = fmt.Sprintf(":%v", val[0])
		default:
			val = data[0:]
		}
		encoder.WriteField(hpack.HeaderField{Name: strings.TrimSpace(val[0]), Value: strings.TrimSpace(val[1])})
	}
	return hbuf.Bytes()
}

// Takes in the url and returns the host + port of the url.
//
//	e.g. "www.google.com:443"
func CheckAddr(url *url.URL) string {
	switch url.Scheme {
	case "https":
		return url.Host + ":443"
	default:
		return url.Host + ":80"
	}
}

// This returns the default config variables.
// header order, chrome like headers and protocols.
func GetDefaultConfig() Config {
	return Config{
		HeaderOrder: []string{
			":authority",
			":method",
			":path",
			":scheme",
			"accept",
			"accept-encoding",
			"accept-language",
			"cache-control",
			"content-length",
			"content-type",
			"cookie",
			"origin",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-ch-ua-platform-version",
			"upgrade-insecure-requests",
			"user-agent",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
		},
		Headers: map[string]string{
			"cache-control":             "max-age=0",
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"sec-fetch-site":            "same-origin",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"sec-fetch-dest":            "document",
			"sec-ch-ua":                 `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`,
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        `"Windows"`,
			"accept-language":           "en-US,en;q=0.9",
		},
		Protocols: []string{"h2", "h1", "http/1.1"},
	}
}

//////////////////////////////////////////////////////////////////////////////////////////

// THIS IS A EXPERIMENTAL CF BYPASS METHOD, USE IF YOU WANT!

//////////////////////////////////////////////////////////////////////////////////////////

type Request struct {
	Port, Url, Host string
	Request         *http.Request
	Data            []byte
	dataNop         io.ReadCloser
	Conn            *tls.UConn
}

func (R *Request) Commit(Redirect bool) (*tls.UConn, *http.Response, error) {

	var request *http.Request = R.Request.Clone(R.Request.Context())
	request.Body = io.NopCloser(bytes.NewBuffer([]byte(R.Data)))
	request.Write(R.Conn)
	read := bufio.NewReader(R.Conn)
	resp, err := http.ReadResponse(read, request)

	if err == nil && Redirect {
		switch resp.StatusCode {
		case 300, 301, 302, 303, 307, 308:
			if res, err := url.Parse(R.Url); err == nil {
				R.ChangeURL(res.Scheme + "://" + res.Host + resp.Header.Get("Location"))
				R.Request.Method = "GET"
				return R.Commit(Redirect)
			}
		}
	}

	return R.Conn, resp, err
}

func (R *Request) ChangeURL(uri string) {
	if u, err := url.Parse(uri); err == nil {
		R.Url = u.String()
		R.Request.URL = u
	}
}

func BuildRequest(URL, method, data string, headers map[string]string, P *ProxyAuth, cookies ...*http.Cookie) (Request, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return Request{}, err
	}

	var port = ":443"
	if strings.EqualFold(u.Scheme, "http") {
		port = ":80"
	}

	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return Request{}, err
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	var conn *tls.UConn
	if P != nil && P.IP != "" && P.Port != "" {
		if conn_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", P.IP, P.Port)); err == nil {
			Conn_URL := CheckAddr(req.URL)
			if P.User != "" && P.Password != "" {
				conn_.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Authorization: Basic %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\nConnection: keep-alive\r\n\r\n", Conn_URL, Conn_URL, base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", P.User, P.Password))))))
			} else {
				conn_.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\nConnection: keep-alive\r\n\r\n", Conn_URL, Conn_URL)))
			}
			var junk = make([]byte, 4096)
			conn_.Read(junk)
			switch Status := string(junk); Status[9:12] {
			case "200":
				conn = tls.UClient(conn_, &tls.Config{
					ServerName: req.Host,
				}, tls.HelloChrome_112, true, true)
			}
		}
	} else {
		if conn_, err := net.Dial("tcp", u.Host+port); err == nil {
			conn = tls.UClient(conn_, &tls.Config{
				ServerName: req.Host,
			}, tls.HelloChrome_112, true, true)
		}
	}

	return Request{
		Port: port, Url: URL, Host: u.Host,
		Request: req,
		Data:    []byte(data),
		Conn:    conn,
		//dataNop: io.NopCloser(bytes.NewBuffer([]byte(data))),
	}, nil

}

func (R *Request) ChangeProxyConn(P *ProxyAuth) {
	var conn *tls.UConn
	if P != nil && P.IP != "" && P.Port != "" {
		if conn_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", P.IP, P.Port)); err == nil {
			Conn_URL := CheckAddr(R.Request.URL)
			if P.User != "" && P.Password != "" {
				conn_.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Authorization: Basic %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\nConnection: keep-alive\r\n\r\n", Conn_URL, Conn_URL, base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", P.User, P.Password))))))
			} else {
				conn_.Write([]byte(fmt.Sprintf("CONNECT %v HTTP/1.1\r\nHost: %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: MCSN/1.1\r\nConnection: keep-alive\r\n\r\n", Conn_URL, Conn_URL)))
			}
			var junk = make([]byte, 4096)
			conn_.Read(junk)
			switch Status := string(junk); Status[9:12] {
			case "200":
				conn = tls.UClient(conn_, &tls.Config{
					ServerName: R.Host,
				}, tls.HelloChrome_112, true, true)
			}
		}
	} else {
		if conn_, err := net.Dial("tcp", R.Host+R.Port); err == nil {
			conn = tls.UClient(conn_, &tls.Config{
				ServerName: R.Host,
			}, tls.HelloChrome_112, true, true)
		}
	}
	R.Conn = conn
}
