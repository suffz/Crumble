package h2

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	tls "github.com/refraction-networking/utls"
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
	}, config.BuildID)

	if config.SaveCookies {
		if Data.Client.Cookies == nil || len(Data.Client.Cookies) == 0 {
			Data.Client.Cookies = make(map[string][]hpack.HeaderField)
		}
	}

	fmt.Fprintf(tlsConn, http2.ClientPreface)

	if err = tlsConn.Handshake(); err != nil {
		return err
	}

	if config.DataBodyMaxLength == 0 {
		return errors.New("error: [DataBodyMaxLength] cannot be 0, suggested value to be above 130000")
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
		cookie_name := strings.Split(val.Value, "=")
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
			Config.Headers, err = hpack.NewDecoder(100000, nil).DecodeFull(f.HeaderBlockFragment())
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
			"cache-control":              "max-age=0",
			"upgrade-insecure-requests":  "1",
			"user-agent":                 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			"accept":                     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"sec-fetch-site":             "same-origin",
			"sec-fetch-mode":             "navigate",
			"sec-fetch-user":             "?1",
			"sec-fetch-dest":             "document",
			"sec-ch-ua":                  `"Google Chrome";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`,
			"sec-ch-ua-mobile":           "?0",
			"sec-ch-ua-platform":         "\\\"Windows\\",
			"sec-ch-ua-platform-version": "14.0.0",
			"accept-language":            "en-US,en;q=0.9",
		},
		Protocols: []string{"h2", "h1", "http/1.1"},
	}
}
