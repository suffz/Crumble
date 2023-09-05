package followbot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"main/packages/utils"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	cu "main/packages/Chrome"

	"github.com/PuerkitoBio/goquery"
	html "github.com/antchfx/htmlquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func X_Forwarded_For() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%v.%v.%v.%v", rand.Intn(300-256)+256, rand.Intn(255), rand.Intn(255), rand.Intn(9))
}

func Get_CF_Clearance() (status int64, cookies []*network.Cookie) {
	if ctx, cancel, err := cu.New(cu.NewConfig(
		cu.WithChromeFlags(append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("headless", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("window-size", "800,600"),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
		)...),
		cu.WithTimeout(time.Second*60),
	)); err == nil {
		defer cancel()
		if err := chromedp.Run(ctx, Perform_chrome(ctx, &status, &cookies)); err != nil {
			if runtime.GOOS != "windows" {

			} else {
				fmt.Println("Please install chrome..")
				os.Exit(0)
			}
		}
	}
	return
}

func Perform_chrome(c context.Context, st *int64, cookies *[]*network.Cookie) chromedp.Tasks {
	chromedp.ListenTarget(c, func(ev interface{}) {
		switch r := ev.(type) {
		case *network.EventResponseReceived:
			resp := r.Response
			if resp.URL == "https://namemc.com/login" {
				*st = resp.Status
			}
		}
	})
	return chromedp.Tasks{
		network.Enable(),
		emulation.SetTouchEmulationEnabled(true),
		emulation.SetEmitTouchEventsForMouse(true),
		chromedp.Navigate("https://namemc.com/login"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			go func(ctx context.Context) {
				var Cont bool
				for !Cont {
					go func() {
						err := chromedp.MouseClickXY(68, 294, chromedp.ButtonLeft).Do(ctx)
						if err != nil {
							Cont = true
						}
					}()
					time.Sleep(1 * time.Second)
				}
			}(ctx)
			return nil
		}),
		chromedp.WaitVisible("#search-box"),
		chromedp.ActionFunc(func(ctx context.Context) error { *cookies, _ = network.GetCookies().Do(ctx); return nil }),
	}
}

func Swapprofile(Profile string) {
	req, _ := http.NewRequest("GET", Profile, nil)
	req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.AddCookie(&http.Cookie{Name: "referrer", Value: Profile})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("Referer", Profile)
	http.DefaultClient.Do(req)
}

func Claim_NAMEMC(url string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.AddCookie(&http.Cookie{Name: "referrer", Value: url})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("Referer", url)
	http.DefaultClient.Do(req)
}

func (Data *Target) Send() {
	req, _ := http.NewRequest("POST", Data.URL, bytes.NewBuffer([]byte(fmt.Sprintf(`profile=%v&task=follow`, Data.UUID))))
	req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.AddCookie(&http.Cookie{Name: "referrer", Value: Data.URL})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("Referer", Data.URL)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Forwarded-For", X_Forwarded_For())
	if resp, err := http.DefaultClient.Do(req); err == nil {
		b, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(b), "Following") {
			fmt.Println(utils.Logo(fmt.Sprintf("<%v> Followed %v [%v]\n", time.Now().Format("05.000"), Data.URL, strings.Split(strings.Split(string(b), `<span class="" translate="no">`)[1], `</span`)[0])))
		}
	}
}

func GetProfiles() {
	req, _ := http.NewRequest("GET", "https://namemc.com/", nil)
	req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: Session.Key,
	})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, _ := http.DefaultClient.Do(req)
	if resp.Status == resp.Status {
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), `<div style="max-height:380px;overflow-y:auto;">`) && strings.Contains(string(body), `/my-profile/switch?`) {
			Session.DisplayName = strings.Split(strings.Split(string(body), `<span class="" translate="no">`)[1], `</span`)[0]
			for _, acc := range strings.Split(strings.Split(strings.Split(string(body), `<div style="max-height:380px;overflow-y:auto;">`)[1], `</div>`)[0], `<a class="dropdown-item"`) {
				if strings.Contains(acc, "src=") {
					Session.Accounts = append(Session.Accounts, UUIDS{
						Name:    strings.Split(strings.Split(acc, `scale=4">`)[1], `</a>`)[0],
						IconPNG: strings.Split(strings.Split(acc, `src="`)[1], `"`)[0],
						URLPath: "https://namemc.com" + strings.Split(strings.Split(acc, `href="`)[1], `"`)[0],
					})
				}
			}
		}
	}
}

func GetInfo(body []byte) (Req NameRequest) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if d, err := html.Parse(strings.NewReader(string(body))); err == nil {
		if b, ok := doc.Find(`#availability-time`).Attr("datetime"); ok {
			if e, ok := doc.Find(`#availability-time2`).Attr("datetime"); ok {
				if t1, err := time.Parse(time.RFC3339, b); err == nil {
					if t2, err := time.Parse(time.RFC3339, e); err == nil {
						if status, searches := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); status.FirstChild != nil && searches.FirstChild != nil {
							var h string
							if head := html.FindOne(d, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
								h = head.Attr[1].Val
							}
							Req = NameRequest{
								Status:     status.FirstChild.Data,
								Searches:   strings.Split(searches.FirstChild.Data, " / month")[0],
								Start:      &t1,
								End:        &t2,
								Start_Unix: t1.Unix(),
								End_Unix:   t2.Unix(),
								HeadURL:    h,
							}
						}
					}
				}
			}
		} else {
			if status, searches := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); status.FirstChild != nil && searches.FirstChild != nil {
				var HeadUrl string
				if head := html.FindOne(d, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
					HeadUrl = head.Attr[1].Val
				}
				body := ""
				if HeadUrl != "" {
					body = fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", strings.Split(strings.Split(HeadUrl, "?id=")[1], "&")[0])
				}
				Req = NameRequest{
					Status:   status.FirstChild.Data,
					Searches: strings.Split(searches.FirstChild.Data, " / month")[0],
					HeadURL:  HeadUrl,
					BodyUrl:  body,
				}
			}
		}
	}
	return
}

func GetAllSkinsOnPage(html string) (Return []NamemcSkins) {
	fmt.Println(html)
	Body := strings.Split(strings.Split(html, `<div class="row gx-2 justify-content-center">`)[1], `<nav>`)[0]
	Resp := Body[:strings.LastIndex(Body, "</div>")]
	for _, packs := range strings.Split(Resp, `<a href="/skin/`) {
		Content := strings.Split(packs, "</a>")[0]
		if strings.Contains(Content, "card-header") {
			ID := strings.Split(Content, `"`)[0]
			var emoji string
			Data := strings.Split(strings.Split(Content, `<span`)[1], ">")[1]
			if strings.Contains(Data, "img") {
				emoji = strings.Split(strings.Split(Data, `alt="`)[1], `"`)[0]
				Data = strings.Split(Data, "<img")[0]
			} else {
				Data = strings.Split(Data, `</span`)[0]
			}
			Name := Data
			SkinNum := "#" + strings.Split(strings.Split(Content, `normal-sm">#`)[1], "</div>")[0]
			Stars := strings.Split(Content, "★<")[0]
			Stars = Stars[strings.LastIndex(Stars, ">")+1:] + "★"
			base_time := strings.Split(strings.Split(Content, "★<")[1], `normal-sm">`)[1]
			time := strings.Split(base_time, "<small>")[0] + strings.Split(strings.Split(base_time, "<small>")[1], "</small>")[0]
			Return = append(Return, NamemcSkins{
				Emoji:          emoji,
				Time:           time,
				NamemcUsername: Name,
				Stars:          Stars,
				Number:         SkinNum,
				HeadURL:        fmt.Sprintf(`https://s.namemc.com/2d/skin/face.png?id=%v&scale=4`, ID),
				BodyURL:        fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", ID),
				SkinDownload:   fmt.Sprintf(`https://s.namemc.com/i/%v.png`, ID),
			})
		}
	}
	return
}

func SendFollow(name string) {
	for _, acc := range Session.Accounts {
		if strings.EqualFold(acc.Name, utils.Con.NMC.Display) {
			Swapprofile(acc.URLPath)
			break
		}
	}
	req, _ := http.NewRequest("GET", "https://namemc.com/profile/"+name, nil)
	req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	var redirect string
	Client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { redirect = req.URL.String(); return nil }}
	resp, _ := Client.Do(req)
	if body, err := io.ReadAll(resp.Body); err == nil && resp.StatusCode == 200 {
		(&Target{
			UUID: strings.Split(strings.Split(string(body), `order-lg-2 col-lg" style="font-size: 90%"><samp>`)[1], `</samp></div>`)[0],
			URL:  redirect,
		}).Send()
	}
}
