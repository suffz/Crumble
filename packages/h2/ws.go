package h2

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"

	tls "github.com/bogdanfinn/utls"
	"github.com/gobwas/ws"
	"github.com/gorilla/websocket"
)

type WebsocketOptions struct {
	URL        string
	ServerName string
	Origin     string
	Host       string
	Extensions string
	UserAgent  string
}

func (Info *WebsocketOptions) Websocket() (*tls.UConn, *http.Response, error) {

	k, err := generateKey()
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.NewRequest("GET", Info.URL, nil)
	if err != nil {
		return nil, nil, err
	}

	addmaddHeaders(resp, map[string]string{
		"Sec-WebSocket-Key":        k,
		"Origin":                   Info.Origin,
		"Upgrade":                  "websocket",
		"Connection":               "Upgrade",
		"Host":                     Info.Host,
		"Sec-WebSocket-Extensions": Info.Extensions,
		"Sec-WebSocket-Version":    "13",
		"User-Agent":               Info.UserAgent,
	})

	if conn, err := net.Dial("tcp", Info.ServerName+":443"); err == nil {
		tlsC := tls.UClient(conn, &tls.Config{
			ServerName: Info.ServerName,
		}, tls.HelloChrome_120_PQ, true, true)
		if websocket.IsWebSocketUpgrade(resp) {
			resp.Write(tlsC)
			resp, err := http.ReadResponse(bufio.NewReader(tlsC), resp)
			return tlsC, resp, err
		} else {
			return nil, nil, errors.New("Request doesnt look like a websocket upgrade.")
		}
	} else {
		return nil, nil, err
	}
}

func addmaddHeaders(resp *http.Request, Headers map[string]string) {
	for key, value := range Headers {
		resp.Header.Add(key, value)
	}
}

func generateKey() (string, error) {
	// 1. 16-byte value
	p := make([]byte, 16)

	// 2. Randomly selected
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}

	// 3. Base64-encoded
	return base64.StdEncoding.EncodeToString(p), nil
}

func SendMessage(conn *tls.UConn, messages []string) error {
	for _, msgs := range messages {
		if err := ws.WriteFrame(conn, ws.MaskFrame(ws.NewTextFrame([]byte(msgs)))); err != nil {
			return err
		}
	}
	return nil
}
