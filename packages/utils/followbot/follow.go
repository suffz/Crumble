package followbot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	cu "main/packages/Chrome"
	"main/packages/h2"
	"main/packages/utils"

	"github.com/PuerkitoBio/goquery"
	html "github.com/antchfx/htmlquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/fogleman/gg"
)

var CommonHeaders = map[string]string{
	"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"accept":                    "application/json, text/javascript, text/html, application/xhtml+xml, application/xml;q=0.9, image/avif, image/webp, image/apng, */*;q=0.8, application/signed-exchange;v=b3",
	"cache-control":             "max-age=0",
	"upgrade-insecure-requests": "1",
	"sec-fetch-site":            "none",
	"sec-fetch-mode":            "navigate",
	"sec-fetch-user":            "?1",
	"sec-fetch-dest":            "document",
	"sec-ch-ua":                 `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`,
	"sec-ch-ua-mobile":          "?0",
	"sec-ch-ua-platform":        `"Windows"`,
	"accept-language":           "en-US,en;q=0.9",
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
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"),
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
	Headers := CommonHeaders
	Headers["referer"] = "https://namemc.com/"
	request, _ := h2.BuildRequest(Profile, "GET", "", Headers,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key}, &http.Cookie{Name: "referrer", Value: "https://namemc.com/"})
	request.Commit()
}

func Claim_NAMEMC(url string) {
	Headers := CommonHeaders
	Headers["origin"] = "https://namemc.com"
	Headers["referer"] = url
	request, _ := h2.BuildRequest(url, "GET", "", Headers,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key}, &http.Cookie{Name: "referrer", Value: url})
	request.Commit()
}

func (Data *Target) Send(precheck bool) bool {

	Headers := CommonHeaders
	Headers["content-type"] = "application/x-www-form-urlencoded"
	Headers["origin"] = "https://namemc.com"
	Headers["referer"] = Data.URL

	request_get, _ := h2.BuildRequest(Data.URL, "GET", "", CommonHeaders, nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})
	request_post, _ := h2.BuildRequest(Data.URL, "POST", fmt.Sprintf(`profile=%v&task=follow`, Data.UUID), Headers, nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key}, &http.Cookie{Name: "referrer", Value: Data.URL})

	var cont bool
	var nn string
	if precheck {
		if _, resp, err := request_get.Commit(); err == nil {
			if resp.Body != nil {
				switch resp.StatusCode {
				case 200:
					if resp.Body != nil {
						ff := resp.Body
						cont = strings.Contains(string(ff), `id="followingMenuButton"`)
						nn = strings.Split(strings.Split(string(ff), `<span class="" translate="no">`)[1], `</span`)[0]
					}
				case 429:
					fmt.Print(utils.Logo("Rate limited, sleeping for 5 seconds.. " + Data.URL + "\n"))
					time.Sleep(5 * time.Second)
					Data.Send(precheck)
				}
			}
		}
	}

	if !cont {
		if _, resp, err := h2.RClient().Redirect(&request_post); err == nil {
			if resp.StatusCode == 429 {
				fmt.Print(utils.Logo("Rate limited, sleeping for 900 seconds.. " + Data.URL + "\n"))
				time.Sleep(900 * time.Second)
				Data.Send(precheck)
			} else {
				switch resp.StatusCode {
				case 429:
					fmt.Print(utils.Logo("Rate limited, sleeping for 900 seconds.. " + Data.URL + "\n"))
					time.Sleep(900 * time.Second)
					Data.Send(precheck)
				default:
					resp = request_post.Redirects
					switch resp.StatusCode {
					case 429:
						fmt.Print(utils.Logo("Rate limited, sleeping for 900 seconds.. " + Data.URL + "\n"))
						time.Sleep(900 * time.Second)
						Data.Send(precheck)
					default:
						if resp.Body != nil {
							b := resp.Body
							if strings.Contains(string(b), `id="followingMenuButton"`) {
								fmt.Print(utils.Logo(fmt.Sprintf("<%v> Followed %v [%v]\n", time.Now().Format("05.000"), Data.URL, strings.Split(strings.Split(string(b), `<span class="" translate="no">`)[1], `</span`)[0])))
								return true
							}
						}
					}
				}
			}
		}
	} else {
		fmt.Print(utils.Logo(fmt.Sprintf("<%v> Already following %v [%v]\n", time.Now().Format("05.000"), Data.URL, nn)))
	}

	return false
}

func GetProfiles() {
	request, _ := h2.BuildRequest("https://namemc.com/", "GET", "", CommonHeaders, nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})
	if _, resp, err := request.Commit(); err == nil {
		body := resp.Body
		if strings.Contains(string(body), `/my-profile/switch?`) {
			Session.DisplayName = strings.Split(strings.Split(string(body), `<span class="" translate="no">`)[1], `</span`)[0]
			for _, acc := range strings.Split(strings.Split(strings.Split(string(body), `<div style="max-height:380px;overflow-y:auto;">`)[1], `</div>`)[0], `<a class="dropdown-item"`) {
				if strings.Contains(acc, "src=") {
					Session.Accounts = append(Session.Accounts, UUIDS{
						Name:    strings.Split(strings.Split(strings.Split(acc, `scale=4">`)[1], ` <span`)[0], "</a>")[0],
						IconPNG: strings.ReplaceAll(strings.Split(strings.Split(acc, `src="`)[1], `"`)[0], "amp;", ""),
						URLPath: strings.ReplaceAll("https://namemc.com"+strings.Split(strings.Split(acc, `href="`)[1], `"`)[0], "amp;", ""),
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

func SendFollow(name string) bool {
	for _, acc := range Session.Accounts {
		if strings.EqualFold(acc.Name, utils.Con.NMC.Display) {
			Swapprofile(acc.URLPath)
			break
		}
	}
	request_get, _ := h2.BuildRequest("https://namemc.com/profile/"+name, "GET", "", CommonHeaders, nil, ReturnCF())
	if _, resp, err := h2.RClient().Redirect(&request_get); err == nil {
		resp = request_get.Redirects
		body := resp.Body
		Headers := CommonHeaders
		Headers["content-type"] = "application/x-www-form-urlencoded"
		Headers["origin"] = "https://namemc.com"
		Headers["referer"] = request_get.Url
		request_post, _ := h2.BuildRequest(request_get.Url, "POST", fmt.Sprintf(`profile=%v&task=follow`, strings.Split(strings.Split(string(body), `order-lg-2 col-lg" style="font-size: 90%"><samp>`)[1], `</samp></div>`)[0]), Headers, nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key}, &http.Cookie{Name: "referrer", Value: request_get.Url})
		if _, resp, err := request_post.Commit(); err == nil {
			return resp.StatusCode == 302
		}

	}
	return false
}

func ParseAndPerformThreeReqs() []Resps {
	var url_char *url.URL = &url.URL{
		Host:       "namemc.com",
		Scheme:     "https",
		Path:       "minecraft-names",
		RawQuery:   "sort=asc&length_op=eq&length=3&lang=&searches=0",
		ForceQuery: true,
	}

	request, _ := h2.BuildRequest(url_char.String(), "GET", "", CommonHeaders,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})

	var all int
	var ThreeCache []Resps

	for {
		var ClientInfo []Resps
		if _, resp, err := request.Commit(); err == nil {
			body := resp.Body
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
				request.ChangeURL(fmt.Sprintf("https://namemc.com%v", html.FindOne(doc, "/html/body/main/div/div[5]/nav/ul/li[4]/a").Attr[1].Val))
				if len(ClientInfo) == 1 {
					break
				} else {
					ThreeCache = append(ThreeCache, ClientInfo...)
				}
			}

		} else {
			time.Sleep(15 * time.Second)
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

	request, _ := h2.BuildRequest("https://namemc.com/search?q="+name, "GET", "", CommonHeaders,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})

	if _, resp, err := request.Commit(); err == nil {
		body := resp.Body
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
	return 0, 0, "", ""
}

func GetHeadUrl(name string) (string, string) {
	request, _ := h2.BuildRequest("https://namemc.com/search?q="+name, "GET", "", CommonHeaders,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})
	if _, resp, err := request.Commit(); err == nil {
		body, _ := html.Parse(bytes.NewBuffer(resp.Body))
		if head := html.FindOne(body, `/html/body/main/div/div/div[3]/div[1]/div/div[2]/img`); head != nil && head.Attr != nil {
			return head.Attr[1].Val, fmt.Sprintf("https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200", strings.Split(strings.Split(head.Attr[1].Val, "?id=")[1], "&")[0])
		}

	}
	return "", ""
}

func Info(name string) (namemc NameRequest) {
	request, _ := h2.BuildRequest("https://namemc.com/search?q="+name, "GET", "", CommonHeaders,
		nil, ReturnCF(), &http.Cookie{Name: "session_id", Value: Session.Key})
	if _, resp, err := request.Commit(); err == nil {
		if resp.Body != nil {
			namemc = GetInfo(resp.Body)
		}
	}
	return
}

func GenImages(body string) []byte {
	art := gg.NewContext(360, 120)
	var test []Images
	c, d, b := 0, 0, 9
	var wg sync.WaitGroup
	for i, skins := range GetAllSkins(strings.Split(strings.Split(body, `<div style="width: 324px; margin: auto; text-align: center;">`)[1], "</div>")[0]) {
		wg.Add(1)
		go func(skins Skins, i int) {
			if img := GetSkin(skins.Head[:len(skins.Head)-1]+"5", i); !reflect.DeepEqual(img, Images{}) {
				test = append(test, img)
			} else {
				for {
					img = GetSkin(skins.Head[:len(skins.Head)-1]+"5", i)
					if !reflect.DeepEqual(img, Images{}) {
						break
					}
				}
			}
			wg.Done()
		}(skins, i)
	}
	wg.Wait()
	sort.Slice(test, func(i, j int) bool {
		return test[i].Row < test[j].Row
	})
	for i, m := range test {
		if i == b {
			b = b + 9
			c = 0
			d = d + m.Image.Bounds().Dy()
		}
		art.DrawImage(m.Image, c, d)
		c = c + m.Image.Bounds().Dx()
	}
	var buf bytes.Buffer
	art.EncodePNG(&buf)
	return buf.Bytes()
}

func GetAllSkins(body string) (N []Skins) {
	for _, T := range strings.Split(body, "<a") {
		switch true {
		case strings.Contains(T, ` href="/skin/`):
			SID := strings.Split(strings.Split(T, `/skin/`)[1], `">`)[0]
			N = append(N, Skins{
				ID:          SID,
				DownloadURL: "https://s.namemc.com/i/" + SID + ".png",
				URL:         "https://namemc.com/skin/" + SID,
				ChangedAt:   strings.Split(strings.Split(T, `title=`)[1][1:], `"`)[0],
				Head:        "https://s.namemc.com/2d/skin/face.png?id=" + SID + "&scale=4",
				Body:        "https://s.namemc.com/3d/skin/body.png?id=" + SID + "&model=classic&width=150&height=200",
			})
		}
	}
	return
}

func GetSkin(url string, i int) Images {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			if m, _, err := image.Decode(resp.Body); err == nil && !isFullyTransparentPng(m) {
				return Images{
					Row:   i,
					Image: m,
					Url:   url,
				}
			} else {
				return GetSkin(url, i)
			}
		} else {
			return GetSkin(url, i)
		}
	}
	return Images{}
}

func isFullyTransparentPng(img image.Image) bool {
	for x := img.Bounds().Min.X; x < img.Bounds().Dx(); x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Dy(); y++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != 0 {
				return false
			}
		}
	}
	return true
}
