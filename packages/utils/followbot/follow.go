package followbot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	cu "main/packages/Chrome"
	"main/packages/utils"

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

func ReturnCF() *http.Cookie {
	return &http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]}
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
		chromedp.Run(ctx, Perform_chrome(ctx, &status, &cookies))
		cancel()
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
	req.AddCookie(ReturnCF())
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.AddCookie(&http.Cookie{Name: "referrer", Value: Profile})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("Referer", Profile)
	http.DefaultClient.Do(req)
}

func Claim_NAMEMC(url string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(ReturnCF())
	req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req.AddCookie(&http.Cookie{Name: "referrer", Value: url})
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("Referer", url)
	http.DefaultClient.Do(req)
}

func (Data *Target) Send() bool {

	// precheck if acc is alr following.
	req2, _ := http.NewRequest("GET", Data.URL, nil)
	req2.AddCookie(ReturnCF())
	req2.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
	req2.Header.Add("X-Forwarded-For", X_Forwarded_For())
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	if bod, err := http.DefaultClient.Do(req2); err == nil && bod.StatusCode == 200 {
		ff, _ := io.ReadAll(bod.Body)
		if !strings.Contains(string(ff), `id="followingMenuButton"`) {
			req, _ := http.NewRequest("POST", Data.URL, bytes.NewBuffer([]byte(fmt.Sprintf(`profile=%v&task=follow`, Data.UUID))))
			req.AddCookie(ReturnCF())
			req.AddCookie(&http.Cookie{Name: "session_id", Value: Session.Key})
			req.AddCookie(&http.Cookie{Name: "referrer", Value: Data.URL})
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
			req.Header.Add("Origin", "https://namemc.com")
			req.Header.Add("Referer", Data.URL)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("X-Forwarded-For", X_Forwarded_For())
			if resp, err := http.DefaultClient.Do(req); err == nil {
				b, _ := io.ReadAll(resp.Body)
				if strings.Contains(string(b), `id="followingMenuButton"`) {
					fmt.Print(utils.Logo(fmt.Sprintf("<%v> Followed %v [%v]\n", time.Now().Format("05.000"), Data.URL, strings.Split(strings.Split(string(b), `<span class="" translate="no">`)[1], `</span`)[0])))
					return true
				}
			}
		}
	}
	return false
}

func GetProfiles() {
	req, _ := http.NewRequest("GET", "https://namemc.com/", nil)
	req.AddCookie(ReturnCF())
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
							body := ""
							if head := html.FindOne(d, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
								h = head.Attr[1].Val
							}
							if h != "" {
								body = fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", strings.Split(strings.Split(h, "?id=")[1], "&")[0])
							}
							Req = NameRequest{
								Status:     status.FirstChild.Data,
								Searches:   strings.Split(searches.FirstChild.Data, " / month")[0],
								Start:      &t1,
								End:        &t2,
								Start_Unix: t1.Unix(),
								End_Unix:   t2.Unix(),
								HeadURL:    h,
								BodyUrl:    body,
							}
						}
					}
				}
			}
		} else {
			if status, searches := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); status != nil && searches != nil && status.FirstChild != nil && searches.FirstChild != nil {
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
	req.AddCookie(ReturnCF())
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

func ParseAndPerformThreeReqs() []Resps {
	var url_char *url.URL = &url.URL{
		Host:       "namemc.com",
		Scheme:     "https",
		Path:       "minecraft-names",
		RawQuery:   "sort=asc&length_op=eq&length=3&lang=&searches=0",
		ForceQuery: true,
	}

	var all int
	var ThreeCache []Resps

	for {

		var (
			ClientInfo []Resps
		)

		req, _ := http.NewRequest("GET", url_char.String(), nil)
		req.AddCookie(ReturnCF())
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		req.Header.Add("Origin", "https://namemc.com")
		req.Header.Add("X-Forwarded-For", X_Forwarded_For())
		req.Header.Add("Referer", "https://namemc.com/")

		if resp, err := http.DefaultClient.Do(req); err == nil {
			if resp.StatusCode == 200 {
				if body, err := io.ReadAll(resp.Body); err == nil {
					if doc, err := html.Parse(strings.NewReader(string(body))); err == nil {
						for i := 1; i <= 119; i += 2 {
							var Searches string
							Search := html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td[7]", i))
							if Search == nil {
								Searches = "0"
							} else {
								if Search.FirstChild != nil {
									Searches = Search.FirstChild.Data
								} else {
									Searches = "0"
								}
							}

							end := html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td/span", i))
							if end != nil {
								st := end.Attr[1].Val
								en := end.Attr[2].Val
								t1, _ := time.Parse(time.RFC3339, st)
								t2, _ := time.Parse(time.RFC3339, en)
								ClientInfo = append(ClientInfo, Resps{
									Name:       strings.ToLower(html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td[1]/a", i)).FirstChild.Data),
									Start:      &t1,
									Start_Unix: t1.Unix(),
									End:        &t2,
									End_Unix:   t2.Unix(),
									Searches:   Searches,
									Index:      all,
								})
								all++
							}
						}
						url, err := url.Parse(fmt.Sprintf("https://namemc.com%v", html.FindOne(doc, "/html/body/main/div/div[5]/nav/ul/li[4]/a").Attr[1].Val))
						if err != nil {
							continue
						}
						url_char = url
						if len(ClientInfo) == 1 {
							break
						} else {
							ThreeCache = append(ThreeCache, ClientInfo...)
						}
					}
				}
			} else {
				time.Sleep(15 * time.Second)
			}
		}
	}

	var removedupes map[string]Resps = make(map[string]Resps)

	for _, name := range ThreeCache {
		removedupes[name.Name] = name
	}

	var New []Resps

	for _, value := range removedupes {
		New = append(New, value)
	}

	sort.Slice(New, func(i, j int) bool {
		return New[i].Index < New[j].Index
	})

	return New
}

func GetDroptimes(name string) (int64, int64, string, string) {

	req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+name, nil)
	req.AddCookie(ReturnCF())
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("X-Forwarded-For", X_Forwarded_For())
	req.Header.Add("Referer", "https://namemc.com/")

	if resp, err := http.DefaultClient.Do(req); err == nil {
		if resp.StatusCode == 200 {
			if body, err := io.ReadAll(resp.Body); err == nil {
				doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
				start, _ := doc.Find("#availability-time").Attr("datetime")
				end, _ := doc.Find("#availability-time2").Attr("datetime")
				s1, _ := time.Parse(time.RFC3339, start)
				s2, _ := time.Parse(time.RFC3339, end)
				d, _ := html.Parse(bytes.NewBuffer(body))
				status, searches := "", ""
				if st, se := html.FindOne(d, `//*[@id="status-bar"]/div/div[1]/div[2]`), html.FindOne(d, `//*[@id="status-bar"]/div/div[2]/div[2]`); st.FirstChild != nil && se.FirstChild != nil {
					status = st.FirstChild.Data
					searches = strings.Split(se.FirstChild.Data, " / month")[0]
				}
				return s1.Unix(), s2.Unix(), status, searches
			}
		}
	}
	return 0, 0, "", ""
}

func GetHeadUrl(name string) (string, string) {
	req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+name, nil)
	req.AddCookie(ReturnCF())
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("X-Forwarded-For", X_Forwarded_For())
	req.Header.Add("Referer", "https://namemc.com/")
	if resp, err := http.DefaultClient.Do(req); err == nil {
		if resp.StatusCode == 200 {
			if body, err := io.ReadAll(resp.Body); err == nil {
				body, _ := html.Parse(bytes.NewBuffer(body))
				if head := html.FindOne(body, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
					return head.Attr[1].Val, fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", strings.Split(strings.Split(head.Attr[1].Val, "?id=")[1], "&")[0])
				}
			}
		}
	}
	return "", ""
}

func Info(name string) (namemc NameRequest) {
	req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+name, nil)
	req.AddCookie(ReturnCF())
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://namemc.com")
	req.Header.Add("X-Forwarded-For", X_Forwarded_For())
	req.Header.Add("Referer", "https://namemc.com/")
	if resp, err := http.DefaultClient.Do(req); err == nil {
		resp_body, _ := io.ReadAll(resp.Body)
		namemc = GetInfo(resp_body)
	}
	return
}
