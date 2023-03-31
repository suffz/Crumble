package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	html "github.com/antchfx/htmlquery"
	"github.com/iskaa02/qalam/gradient"
)

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

func IsAvailable(name string) bool {
	resp, err := http.Get("https://account.mojang.com/available/minecraft/" + name)
	if err == nil {
		return resp.StatusCode == 200
	} else {
		return false
	}
}

func Logo(Data string) string {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(RGB...).
		Mode(gradient.BlendRgb).
		Build()
	return g.Mutline(Data)
}

func GetHeadUrl(name, cookies string) (string, string) {
	req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+name, nil)
	for _, cookies := range strings.Split(cookies, "; ") {
		req.AddCookie(&http.Cookie{
			Name:  strings.Split(cookies, "=")[0],
			Value: strings.ReplaceAll(strings.ReplaceAll(strings.Split(cookies, "=")[1], `"`, ""), "\x00", ""),
		})
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	if resp, err := http.DefaultClient.Do(req); err == nil {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			if d, err := html.Parse(strings.NewReader(string(body))); err == nil {
				if head := html.FindOne(d, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
					return head.Attr[1].Val, fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", strings.Split(strings.Split(head.Attr[1].Val, "?id=")[1], "&")[0])
				} else {
					return "", ""
				}
			}
		}
	} else {
		fmt.Println(err)
	}

	return "", ""
}

func GetDroptimes(name, cookies string) (int64, int64, string, string) {
	req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+name, nil)
	for _, cookies := range strings.Split(cookies, "; ") {
		req.AddCookie(&http.Cookie{
			Name:  strings.Split(cookies, "=")[0],
			Value: strings.ReplaceAll(strings.ReplaceAll(strings.Split(cookies, "=")[1], `"`, ""), "\x00", ""),
		})
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	var st, se string
	if resp, err := http.DefaultClient.Do(req); err == nil {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
			if d, err := html.Parse(strings.NewReader(string(body))); err == nil {
				if status, searches := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); status.FirstChild != nil && searches.FirstChild != nil {
					st = status.FirstChild.Data
					se = strings.Split(searches.FirstChild.Data, " / month")[0]
				}
			}
			if b, ok := doc.Find(`#availability-time`).Attr("datetime"); ok {
				if e, ok := doc.Find(`#availability-time2`).Attr("datetime"); ok {
					if t1, err := time.Parse(time.RFC3339, b); err == nil {
						if t2, err := time.Parse(time.RFC3339, e); err == nil {
							return t1.Unix(), t2.Unix(), st, se
						}
					}
				}
			}
		}
	} else {
		fmt.Println(err)
	}

	return 0, 0, st, se
}
