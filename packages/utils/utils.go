package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/packages/apiGO"
	"main/packages/h2"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/chat"
	"github.com/iskaa02/qalam/gradient"
	"github.com/playwright-community/playwright-go"
)

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
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

func Success() gradient.Gradient {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors("rgb(0,79,0)").
		Mode(gradient.BlendRgb).
		Build()
	return g
}

func Failure() gradient.Gradient {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors("rgb(128,0,21)").
		Mode(gradient.BlendRgb).
		Build()
	return g
}

func Appendinvalids(invalids []string) {
	if _, err := os.Stat("data/invalids.txt"); os.IsNotExist(err) {
		os.Create("data/invalids.txt")
	}
	os.WriteFile("data/invalids.txt", []byte(strings.Join(invalids, "\n")), 0644)
}

func IsChangeable(proxy, bearer string) bool {
	data := strings.Split(proxy, ":")
	var Client http.Client
	switch len(data) {
	case 4:
		Client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", User: url.UserPassword(data[2], data[3]), Host: data[0] + ":" + data[1]})}}
	case 2:
		Client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: data[0] + ":" + data[1]})}}
	default:
		Client = *http.DefaultClient
	}
	req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/namechange", nil)
	req.Header.Add("Authorization", "Bearer "+bearer)
	if resp, err := Client.Do(req); err == nil {
		body, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case 200:
			var Data struct {
				NC bool `json:"nameChangeAllowed"`
			}
			json.Unmarshal(body, &Data)
			return Data.NC
		case 404:
			return true
		}
	}
	return false
}

func Regenerateallaccs() {
	var use_proxy, ug int
	var f bool = true
	var Data []Payload_auth
	for i, bearer := range Con.Bearers {
		if !bearer.NOT_ENTITLED_CHECKED {
			Con.Bearers[i].NOT_ENTITLED_CHECKED = true
			if use_proxy >= len(Proxy.Proxys) && len(Proxy.Proxys) < len(Bearer.Details) {
				break
			}

			if f {
				Data = append(Data, Payload_auth{Proxy: Proxy.Proxys[use_proxy]})
				f = false
				use_proxy++
			}
			if len(Data[ug].Accounts) != 3 {
				Data[ug].Accounts = append(Data[ug].Accounts, bearer.Bearer)
			} else {
				ug++
				Data = append(Data, Payload_auth{Proxy: Proxy.Proxys[use_proxy], Accounts: []string{bearer.Bearer}})
				use_proxy++
			}
		}
	}

	for _, proxy := range Data {
		go func(proxy Payload_auth) {
			for _, account := range proxy.Accounts {
				req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/entitlements/mcstore", nil)
				req.Header.Add("Authorization", "Bearer "+account)
				p := strings.Split(proxy.Proxy, ":")
				var P_url *url.URL
				switch len(p) {
				case 2:
					P_url = &url.URL{
						Scheme: "http",
						Host:   p[0] + ":" + p[1],
					}
				case 4:
					P_url = &url.URL{
						Scheme: "http",
						Host:   p[0] + ":" + p[1],
						User:   url.UserPassword(p[2], p[3]),
					}
				}
				(&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(P_url)}}).Do(req)
			}
		}(proxy)
	}
	Con.SaveConfig()
	Con.LoadState()
}

func ReturnJustString(data []byte, err error) string {
	return string(data)
}

func ReturnAll(bearer string, PS *apiGO.ProxyMS) (Data apiGO.UserINFO, Accounttype string) {
	if PS != nil {
		req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
		req.Header.Add("Authorization", "Bearer "+bearer)
		if resp, err := (&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: PS.IP + ":" + PS.Port, User: url.UserPassword(PS.User, PS.Password)})}}).Do(req); err == nil {
			switch resp.StatusCode {
			case 200:
				json.Unmarshal([]byte(ReturnJustString(io.ReadAll(resp.Body))), &Data)
				Accounttype = "Microsoft"
			case 404:
				Accounttype = "Giftcard"
			}
		}
	} else {
		req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
		req.Header.Add("Authorization", "Bearer "+bearer)
		if resp, err := http.DefaultClient.Do(req); err == nil {
			switch resp.StatusCode {
			case 200:
				json.Unmarshal([]byte(ReturnJustString(io.ReadAll(resp.Body))), &Data)
				Accounttype = "Microsoft"
			case 404:
				Accounttype = "Giftcard"
			}
		}
	}
	return
}

func NotWhiteSpace(str string) bool {
	for _, c := range str {
		if !unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func Namemc_key(bearer string) string {

	var Info_Acc apiGO.UserINFO

	if len(Proxy.Proxys) > 0 {
		ip, port, user, pass := GetProxyStrings(Proxy.CompRand())
		Info_Acc, _ = apiGO.ReturnAll(bearer, &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass})
	} else {
		Info_Acc, _ = apiGO.ReturnAll(bearer, nil)
	}

	client = bot.NewClient()
	client.Auth = bot.Auth{
		Name: Info_Acc.Name,
		UUID: Info_Acc.ID,
		AsTk: bearer,
	}

	var P *basic.Player
	L := basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{
		GameStart: func() error {
			return nil
		},
		Disconnect: func(reason chat.Message) error {
			return nil
		},
		Death: func() error {
			return P.Respawn()
		},
	})
	P = L
	chatHandler = msg.New(client, P, playerlist.New(client), msg.EventsHandler{
		SystemChat: func(c chat.Message, overlay bool) error {
			if Text := c.ClearString(); NotWhiteSpace(Text) {
				if strings.Contains(Text, Info_Acc.Name+" joined the game.") {
					if err := chatHandler.SendMessage("/namemc"); err != nil {
						client.Close()
						return err
					}
				} else if strings.Contains(Text, "https://namemc.com/claim?key=") {
					key := Text
					client.Close()
					return errors.New("got-key:" + key)
				}
			}
			return nil
		},
		PlayerChatMessage: func(msg chat.Message, validated bool) error {
			return nil
		},
		DisguisedChat: func(msg chat.Message) error {
			return nil
		},
	})
	if err := client.JoinServer("blockmania.com"); err == nil {
		for {
			if err := client.HandleGame(); err == nil {
				panic("HandleGame never return nil")
			} else if strings.Contains(err.Error(), "got-key") {
				return strings.Split(err.Error(), "got-key:")[1]
			}
		}
	} else {
		return err.Error()
	}
}

func AddRecoveryEmails(accs_used []apiGO.Info) {
	resp, _ := http.Get("https://www.1secmail.com/api/v1/?action=getDomainList")
	body, _ := io.ReadAll(resp.Body)
	var Domains []string
	json.Unmarshal(body, &Domains)

	if pw, err := playwright.Run(); err == nil {
		if browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Channel:  &[]string{"chrome"}[0],
			Headless: &[]bool{true}[0],
		}); err == nil {
			for _, acc := range accs_used {
				go func(email, password string) {
					if page, err := browser.NewPage(playwright.BrowserNewPageOptions{
						UserAgent: &[]string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"}[0],
					}); err == nil {
						page.SetDefaultTimeout(2500)
						if _, err := page.Goto("https://www.microsoft.com/rpsauth/v1/account/SignIn?ru=https%3A%2F%2Fwww.microsoft.com%2Fen-gb%2F"); err == nil {
							defer page.Close()
							time.Sleep(2 * time.Second)
							page.WaitForSelector("#i0116")
							page.Fill("#i0116", email)
							time.Sleep(2 * time.Second)
							page.WaitForSelector("#idSIButton9")
							page.Click("#idSIButton9")
							time.Sleep(2 * time.Second)
							page.WaitForSelector("#i0118")
							page.Fill("#i0118", password)
							time.Sleep(2 * time.Second)
							page.WaitForSelector("#idSIButton9")
							page.Click("#idSIButton9")

							if c, err := page.WaitForSelector("#iLandingViewAction"); err == nil {
								c.Click()
							}

							page.WaitForSelector("#idBtn_Back")
							page.Click("#idBtn_Back")

							if _, err := page.Goto("https://account.live.com/proofs/manage/additional?mkt=en-US&refd=account.microsoft.com&refp=security"); err == nil {
								if c, err := page.WaitForSelector("#Add_email"); err == nil {
									if err := c.Click(); err != nil {
										fmt.Println(err, email)
										page.Close()
										return
									}
								} else if eval, err := page.WaitForSelector("#iLandingViewTitle"); err == nil {
									if content, err := eval.TextContent(); err == nil {
										if strings.EqualFold(content, "You can't access this site right now") {
											page.Close()
											return
										}
									}
								} else if _, err := page.WaitForSelector("#idDiv_SAOTCS_Proofs > div > div > div > div.table-cell.text-left.content"); err == nil {
									page.Close()
									fmt.Println(Logo(fmt.Sprintf("<%v> %v page required email verification to login..", time.Now().Format("05.000"), email)))
									return
								} else if _, err := page.WaitForSelector("body > pre"); err == nil {
									page.Close()
									return
								}
								email_code := GenerateEmail(Domains)
								if _, err := page.WaitForSelector("#EmailAddress"); err == nil {
									if err := page.Fill("#EmailAddress", email_code); err == nil {
										if err := page.Click("#iNext"); err == nil {
										Exit:
											for i := 0; i < 100; i++ {
												if resp, err := http.Get(fmt.Sprintf("https://www.1secmail.com/api/v1/?action=getMessages&login=%v&domain=%v", strings.Split(email_code, "@")[0], strings.Split(email_code, "@")[1])); err == nil {
													body, _ := io.ReadAll(resp.Body)
													var Body []struct {
														ID      int    `json:"id"`
														From    string `json:"from"`
														Subject string `json:"subject"`
													}
													json.Unmarshal(body, &Body)
													for _, sub := range Body {
														if sub.Subject == "Microsoft account security code" {
															resp, _ := http.Get(fmt.Sprintf("https://www.1secmail.com/api/v1/?action=readMessage&login=%v&domain=%v&id=%v", strings.Split(email_code, "@")[0], strings.Split(email_code, "@")[1], sub.ID))
															body, _ := io.ReadAll(resp.Body)
															var Body struct {
																Body string `json:"body"`
															}
															json.Unmarshal(body, &Body)
															code := regexp.MustCompile("Security code: [0-9]+").FindAllStringSubmatch(string(body), -1)[0][0][15:]
															page.Fill("#iOttText", code)
															page.Click("#iNext")

															fmt.Println(Logo(fmt.Sprintf("<%v> i Added recovery onto "+email+" >> "+email_code, time.Now().Format("05.000"))))
															Con.Recovery = append(Con.Recovery, Succesful{
																Email:     email,
																Recovery:  email_code,
																Code_Used: code,
															})
															Con.SaveConfig()
															Con.LoadState()
															page.Close()
															break Exit
														}
													}
												}
												fmt.Println(Logo(fmt.Sprintf("<%v> Attempt %v/100 to find a email code for %v >> %v", time.Now().Format("05.000"), i, email, email_code)))
												time.Sleep(5 * time.Second)
											}
										}
									}
								} else {
									if Con.Bools.ApplyNewRecoveryToExistingEmails {
										if _, err := page.WaitForSelector("#idA_SAOTCS_LostProofs"); err == nil {
											page.Click("#idA_SAOTCS_LostProofs")
											page.WaitForSelector("#idSubmit_SAOTCS_SendCode")
											page.Click("#idSubmit_SAOTCS_SendCode")
											page.WaitForSelector("#EmailAddress")
											page.Fill("#EmailAddress", email_code)
											page.Click("#iCollectProofAction")
										}
									Exit2:
										for i := 0; i < 100; i++ {
											if resp, err := http.Get(fmt.Sprintf("https://www.1secmail.com/api/v1/?action=getMessages&login=%v&domain=%v", strings.Split(email_code, "@")[0], strings.Split(email_code, "@")[1])); err == nil {
												body, _ := io.ReadAll(resp.Body)
												var Body []struct {
													ID      int    `json:"id"`
													From    string `json:"from"`
													Subject string `json:"subject"`
												}
												json.Unmarshal(body, &Body)
												for _, sub := range Body {
													if sub.Subject == "Microsoft account security info" {
														resp, _ := http.Get(fmt.Sprintf("https://www.1secmail.com/api/v1/?action=readMessage&login=%v&domain=%v&id=%v", strings.Split(email_code, "@")[0], strings.Split(email_code, "@")[1], sub.ID))
														body, _ := io.ReadAll(resp.Body)
														var Body struct {
															Body string `json:"body"`
														}
														json.Unmarshal(body, &Body)
														page.Goto(strings.Split(strings.Split(Body.Body, `href="`)[1], `"`)[0])
														time.Sleep(2 * time.Second)
														page.Fill("#AccountNameInput", email)
														page.Click("#iCollectMembernameAction")
														fmt.Println(Logo(fmt.Sprintf("<%v> i Changed recovery for "+email+" >> "+email_code, time.Now().Format("05.000"))))
														Con.Recovery = append(Con.Recovery, Succesful{
															Email:     email,
															Recovery:  email_code,
															Code_Used: "nil",
														})
														Con.SaveConfig()
														Con.LoadState()
														page.Close()
														break Exit2
													}
												}
											}
											fmt.Println(Logo(fmt.Sprintf("<%v> Attempt %v/100 to find a email code for %v >> %v", time.Now().Format("05.000"), i, email, email_code)))
											time.Sleep(10 * time.Second)
										}
									}
								}
							}
						}
					}
				}(acc.Email, acc.Password)

				if len(browser.Contexts()) == 15 {
					for {
						if len(browser.Contexts()) < 5 {
							break
						}
					}
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			fmt.Println(Logo(err.Error()))
		}
	} else {
		fmt.Println(Logo(err.Error()))
	}
}

func GenerateEmail(domains []string) string {
	rand.Seed(time.Now().UnixMicro())
	return RandStringRunes(64) + "@" + domains[rand.Intn(len(domains))]
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixMicro())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReturnPayload(acc, bearer, name string) string {
	if acc == "Giftcard" {
		var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
		if Con.Bools.UseBypass {
			return fmt.Sprintf("POST https://minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net/minecraft/profile// HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		} else {
			return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		}
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: Bearer " + bearer + "\r\n"
	}
}

func ReturnRequest(acc, bearer, name string, PROXY *h2.ProxyAuth) h2.Request {
	if acc == "Giftcard" {
		var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
		if Con.Bools.UseBypass {
			if req, err := h2.BuildRequest("https://api.minecraftservices.com/minecraft/profile//", "POST", JSON, map[string]string{
				"Content-Length": strconv.Itoa(len(JSON)),
				"Content-Type":   "application/json",
				"Authorization":  "Bearer " + bearer,
			}, PROXY); err == nil {
				return req
			}
		} else {
			if req, err := h2.BuildRequest("https://api.minecraftservices.com/minecraft/profile", "POST", JSON, map[string]string{
				"Content-Length": strconv.Itoa(len(JSON)),
				"Content-Type":   "application/json",
				"Authorization":  "Bearer " + bearer,
			}, PROXY); err == nil {
				return req
			}
		}
	} else {

		if req, err := h2.BuildRequest("https://api.minecraftservices.com/minecraft/profile/name/"+name, "PUT", "", map[string]string{
			"Authorization":  "Bearer " + bearer,
			"Content-Length": "0",
		}, PROXY); err == nil {
			return req
		}
	}
	return h2.Request{}
}
