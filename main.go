package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"main/utils"
	"main/webhook"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	cu "main/Chrome"
	"main/StrCmd"
	"main/apiGO"

	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
)

type Target struct {
	UUID string
	URL  string
	Hits int
}

func TempCalc(interval, accamt int) time.Duration {
	amt := interval / accamt
	if amt < 10 {
		return time.Duration(amt) * time.Millisecond
	}
	return time.Duration((amt)-10) * time.Millisecond
}

var Cookie string

type Namemc struct {
	Key         string
	DisplayName string
	Accounts    []UUIDS
}

type UUIDS struct {
	Name    string
	URLPath string
	IconPNG string
}

var (
	Session Namemc
)

func init() {
	utils.Roots.AppendCertsFromPEM(utils.ProxyByte)
	apiGO.Clear()
	utils.Con.LoadState()

	var reloadifnill bool
	if utils.Con.Settings.GC_ReqAmt == 0 {
		utils.Con.Settings.GC_ReqAmt = 1
		reloadifnill = true
	}

	if utils.Con.Settings.MFA_ReqAmt == 0 {
		utils.Con.Settings.MFA_ReqAmt = 1
		reloadifnill = true
	}

	if reloadifnill {
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}

	for _, rgb := range utils.Con.Gradient {
		utils.RGB = append(utils.RGB, fmt.Sprintf("rgb(%v,%v,%v)", rgb.R, rgb.G, rgb.B))
	}

	if !utils.Con.Bools.DownloadedPW {
		if err := playwright.Install(&playwright.RunOptions{Verbose: true}); err == nil {
			utils.Con.Bools.DownloadedPW = true
			utils.Con.SaveConfig()
			utils.Con.LoadState()
		}
		if utils.Con.NMC.UseNMC {
			if runtime.GOOS != "windows" {
				fmt.Println("Installing google chrome...")
				if _, err := os.Stat("google-chrome-stable_current_amd64.deb"); os.IsNotExist(err) {
					exec.Command("wget", "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb").Run()
					if err := exec.Command("/bin/sh", "-c", "sudo dpkg -i google-chrome-stable_current_amd64.deb").Run(); err != nil {
						exec.Command("/bin/sh", "-c", "sudo apt update").Run()
						exec.Command("/bin/sh", "-c", "sudo apt upgrade").Run()
						exec.Command("/bin/sh", "-c", "sudo dpkg --configure -a").Run()
						exec.Command("/bin/sh", "-c", "sudo apt --fix-broken install -y").Run()
						exec.Command("/bin/sh", "-c", "sudo apt install xvfb").Run()
					}
				}
				exec.Command("/bin/sh", "-c", "Xvfb :1 -screen 0 800x600x24").Run()
				exec.Command("/bin/sh", "-c", "sudo apt-get install -f").Run()
				fmt.Println("Please rerun the program with this command: DISPLAY=:1 go run . *OR* DISPLAY=:1 ./the-name-of-this-executable")
			}
			os.Exit(0)
		}
	}

	fmt.Print(utils.Logo(`                                                    
 ,-----.                         ,--.   ,--.        
'  .--./,--.--.,--.,--.,--,--,--.|  |-. |  | ,---.  
|  |    |  .--'|  ||  ||        || .-. '|  || .-. : 
'  '--'\|  |   '  ''  '|  |  |  || '-' ||  |\   --. 
 '-----''--'    '----' '--'--'--' '---' '--' '----' 
                                                    
`))

	if utils.Con.Bools.FirstUse {
		fmt.Print(utils.Logo("Use proxys for authentication? : [YES/NO] > "))
		var ProxyAuth string
		fmt.Scan(&ProxyAuth)
		utils.Con.Bools.FirstUse = false
		utils.Con.Bools.UseProxyDuringAuth = strings.Contains(strings.ToLower(ProxyAuth), "y")
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}
	if file_name := "accounts.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if file_name := "proxys.txt"; utils.CheckForValidFile(file_name) {
		os.Create(file_name)
	}
	if _, err := os.Stat("skins"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("skins", os.ModePerm)
	}
	if _, err := os.Stat("skinarts"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("skinarts", os.ModePerm)
	}

	utils.Proxy.GetProxys(false, nil)
	utils.Proxy.Setup()
	utils.AuthAccs()
	utils.Regenerateallaccs()
	go utils.CheckAccs()
	var use_proxy, gcamt, mfaamt int

	for _, bearer := range utils.Bearer.Details {
		if use_proxy >= len(utils.Proxy.Proxys) && len(utils.Proxy.Proxys) < len(utils.Bearer.Details) {
			break
		}
		switch bearer.AccountType {
		case "Microsoft":
			utils.Accs["Microsoft"] = append(utils.Accs["Microsoft"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{bearer}})
			utils.Accamt++
			mfaamt++
		case "Giftcard":
			if utils.First_gc {
				utils.Accs["Giftcard"] = []utils.Proxys_Accs{{Proxy: utils.Proxy.Proxys[use_proxy]}}
				utils.First_gc = false
				use_proxy++
			}
			var am int = utils.Con.Settings.AccountsPerGc
			if am == 0 {
				am = 1
			}
			if len(utils.Accs["Giftcard"][utils.Use_gc].Accs) != am {
				utils.Accs["Giftcard"][utils.Use_gc].Accs = append(utils.Accs["Giftcard"][utils.Use_gc].Accs, bearer)
				utils.Accamt++
				gcamt++
			} else {
				utils.Use_gc++
				utils.Accamt++
				gcamt++
				utils.Accs["Giftcard"] = append(utils.Accs["Giftcard"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{bearer}})
				use_proxy++
			}
		}
	}

	if gcamt == 0 {
		gcamt = 1
	}
	if mfaamt == 0 {
		mfaamt = 1
	}

	if utils.Con.NMC.UseNMC {
		if utils.Con.CF.Tokens != "" && time.Now().Unix() < utils.Con.CF.GennedAT {
			Cookie = utils.Con.CF.Tokens
			go func() {
				for {
					time.Sleep(time.Until(time.Unix(utils.Con.CF.GennedAT, 0)))
					if status, cookies := Get_CF_Clearance(); status == 200 {
						for _, cookies := range cookies {
							if cookies.Name == "cf_clearance" {
								utils.Con.CF = utils.CF{
									Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
									GennedAT: time.Now().Add(time.Second * 1800).Unix(),
								}
								Cookie = utils.Con.CF.Tokens
								utils.Con.SaveConfig()
								utils.Con.LoadState()
								break
							}
						}
					}
					time.Sleep(10 * time.Second)
				}
			}()
		} else {
			if status, cookies := Get_CF_Clearance(); status == 200 {
				for _, cookies := range cookies {
					if cookies.Name == "cf_clearance" {
						utils.Con.CF = utils.CF{
							Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
							GennedAT: time.Now().Add(time.Second * 1800).Unix(),
						}
						Cookie = utils.Con.CF.Tokens
						utils.Con.SaveConfig()
						utils.Con.LoadState()
						break
					}
				}
				go func() {
					for {
						time.Sleep(time.Until(time.Unix(utils.Con.CF.GennedAT, 0)))
						if status, cookies := Get_CF_Clearance(); status == 200 {
							for _, cookies := range cookies {
								if cookies.Name == "cf_clearance" {
									utils.Con.CF = utils.CF{
										Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
										GennedAT: time.Now().Add(time.Second * 1800).Unix(),
									}
									Cookie = utils.Con.CF.Tokens
									utils.Con.SaveConfig()
									utils.Con.LoadState()
									break
								}
							}
						}
						time.Sleep(10 * time.Second)
					}
				}()
			}
		}
		if !(time.Now().Unix() < utils.Con.NMC.NamemcLoginData.LastAuthed) {
			if utils.Con.NMC.Key != "" {
				acc := strings.Split(utils.Con.NMC.Key, ":")
				req, _ := http.NewRequest("POST", "https://namemc.com/login", bytes.NewBuffer([]byte(fmt.Sprintf(`email=%v&password=%v`, acc[0], acc[1]))))
				req.AddCookie(&http.Cookie{Name: strings.Split(Cookie, "=")[0], Value: strings.Split(Cookie, "=")[1]})
				req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Add("Origin", "https://namemc.com")
				req.Header.Add("Referer", "https://namemc.com/login")
				Cl := http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						for _, name := range req.Response.Cookies() {
							if name.Name == "session_id" {
								Session.Key = name.Value
							}
						}
						return nil
					},
				}
				r, _ := Cl.Do(req)
				if r.StatusCode == 401 {
					fmt.Println(utils.Logo(fmt.Sprintf("<%v> Email and password invalid for namemc..", time.Now().Format("05.000"))))
				} else {
					if Session.Key != "" {
						utils.Con.NMC.NamemcLoginData = utils.NMC{
							Token:      Session.Key,
							LastAuthed: time.Now().Add(time.Second * 86400).Unix(),
						}
						utils.Con.SaveConfig()
						utils.Con.LoadState()
					}
				}
			}
		} else {
			Session.Key = utils.Con.NMC.NamemcLoginData.Token
		}
		getProfiles()
	}
	fmt.Print(utils.Logo(fmt.Sprintf(`
i Accounts Loaded  > <%v>
i Proxies Loaded   > <%v>
i Proxys in use    > <%v>
i Recoverys Done   > <%v>
i Accounts Details:
 - GC's Per Proxy  > <%v>
 - Req per GC      > <%v>
 - Req per MFA     > <%v>
 - Spread GC       > <%v>
 - Spread MFA      > <%v>
i Namemc Info:
 - Using NameMC    > <%v>
 - Accounts Loaded > <%v>
 - Display Name    > <%v>

`,
		len(utils.Bearer.Details),
		len(utils.Proxy.Proxys),
		use_proxy,
		len(utils.Con.Recovery),
		utils.Con.Settings.AccountsPerGc,
		utils.Con.Settings.GC_ReqAmt,
		utils.Con.Settings.MFA_ReqAmt,
		TempCalc((utils.Con.Settings.SleepAmtPerGc), gcamt),
		TempCalc((utils.Con.Settings.SleepAmtPerMfa), mfaamt),
		utils.Con.NMC.UseNMC,
		len(Session.Accounts),
		Session.DisplayName,
	)))
}

func main() {
	app := StrCmd.App{
		Display:        utils.Logo("@Crumble/root: "),
		Version:        "v1.5.15-CR",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"key": {
				Action: func() {
					var account string
					fmt.Print(utils.Logo("Account [email:password]: "))
					fmt.Scan(&account)
					if details := strings.Split(account, ":"); len(details) >= 2 {
						if info := apiGO.MS_authentication(details[0], details[1], nil); info.Error == "" {
							fmt.Println(utils.Namemc_key(info.Bearer))
						} else {
							fmt.Println(utils.Logo(info.Error))
						}
					}
				},
			},
			"skinart-png": {
				Action: func() {
					var name string
					fmt.Print(utils.Logo("Name of the profile you wanna scrape: "))
					fmt.Scan(&name)
					resp, _ := http.Get("https://namemc.info/data/namemc/skinart/logo/" + name)
					img, _, _ := image.Decode(resp.Body)
					path := "skinarts/" + strings.ReplaceAll(uuid.NewString(), "-", "") + "_" + name + ".png"
					out, _ := os.Create(path)
					png.Encode(out, img)
				},
			},
			"recover": {Action: func() {
				force := StrCmd.Bool("--force")
				use := []apiGO.Info{}
				if force {
					use = utils.Bearer.Details
				} else {
					if _, err := os.Stat("invalids.txt"); !os.IsNotExist(err) {
						body, _ := os.ReadFile("invalids.txt")
						scanner := bufio.NewScanner(bytes.NewBuffer(body))
						var a []string
						for scanner.Scan() {
							a = append(a, scanner.Text())
						}
						for _, acc := range a {
							var Found bool
							for _, accs := range utils.Con.Recovery {
								if strings.EqualFold(accs.Email, strings.Split(acc, ":")[0]) {
									Found = true
									break
								}
							}
							if !Found {
								use = append(use, apiGO.Info{
									Email:    strings.Split(acc, ":")[0],
									Password: strings.Split(acc, ":")[1],
								})
							}
						}
					}
				}
				AddRecoveryEmails(use)
			}, Args: []string{"--force"}},
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					if len(utils.Con.Bearers) == 0 && len(utils.Bearer.Details) == 0 {
						return
					}
					cl, name, Changed, EmailClaimed := false, StrCmd.String("-u"), false, ""
					var start, end int64 = int64(StrCmd.Int("-start")), int64(StrCmd.Int("-end"))
					if utils.Con.NMC.UseNMC {
						start, end, _, _ = utils.GetDroptimes(name)
					} else if !utils.Con.NMC.UseNMC || start == 0 || end == 0 {
						if utils.Con.Settings.AskForUnixPrompt {
							fmt.Println(utils.Logo("Timestamp to Unix: [https://www.epochconverter.com/] (make sure to remove the • on the namemc timestamp!)"))
							fmt.Print(utils.Logo("Use your own unix timestamps [y/n]: "))
							var Use string
							fmt.Scan(&Use)
							if strings.Contains(strings.ToLower(Use), "y") {
								fmt.Print(utils.Logo("Start: "))
								fmt.Scan(&start)
								fmt.Print(utils.Logo("End: "))
								fmt.Scan(&end)
							}
						}
					}
					drop := time.Unix(int64(start), 0)
					for time.Now().Before(drop) {
						fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v                 \r", name, time.Until(drop).Round(time.Second)))))
						time.Sleep(time.Second * 1)
					}
					go func() {
					Exit:
						for {
							if utils.IsAvailable(name) {
								Changed = true
								cl = true
								break Exit
							}
							if start != 0 && end != 0 && time.Now().After(time.Unix(int64(end), 0)) {
								Changed = true
								cl = true
								break Exit
							}
							time.Sleep(10 * time.Second)
						}
					}()
					go func() {
						type Proxys struct {
							Conn     *tls.Conn
							Accounts []apiGO.Info
							Proxy    string
						}
						var Payloads []Proxys
						var wg_a sync.WaitGroup
						for _, Acc := range append(utils.Accs["Giftcard"], utils.Accs["Microsoft"]...) {
							wg_a.Add(1)
							go func(Acc utils.Proxys_Accs) {
								if P, ok := utils.Connect(Acc.Proxy); ok {
									fmt.Println(utils.Logo(fmt.Sprintf("<%v> %v [OK] > %v accs..", time.Now().Format("05.000"), strings.Split(Acc.Proxy, ":")[0], len(Acc.Accs))))
									Payloads = append(Payloads, Proxys{
										Conn:     P,
										Accounts: Acc.Accs,
										Proxy:    Acc.Proxy,
									})
								} else {
									fmt.Println(utils.Logo(fmt.Sprintf("<%v> %v Proxy timed out or couldnt connect.. ", time.Now().Format("05.000"), strings.Split(Acc.Proxy, ":")[0])))
								}
								wg_a.Done()
							}(Acc)
						}

						wg_a.Wait()

						go func() {
							for {
								var wg sync.WaitGroup
								for i, c := range Payloads {
									wg.Add(1)
									go func(i int, c Proxys) {
										if P, ok := utils.Connect(c.Proxy); ok {
											Payloads[i].Conn = P
										}
										wg.Done()
									}(i, c)
								}
								wg.Wait()
								time.Sleep(5 * time.Second)
							}
						}()
						for !cl || !Changed {
							for _, c := range Payloads {
								for _, accs := range c.Accounts {
									if !cl {
										go func(Config apiGO.Info, c Proxys) {
											reqamt := 1
											switch Config.AccountType {
											case "Giftcard":
												reqamt = utils.Con.Settings.GC_ReqAmt
											case "Microsoft":
												reqamt = utils.Con.Settings.MFA_ReqAmt
											}

											if P, ok := utils.Connect(c.Proxy); ok {
												var wg sync.WaitGroup
												for i := 0; i < reqamt; i++ {
													wg.Add(1)
													go func() {
														var status string
														defer wg.Done()
														Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(P, ReturnPayload(Config.AccountType, Config.Bearer, name)), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType}
														switch true {
														case strings.Contains(Req.ResponseDetails.Body, "DUPLICATE"):
															status = "DUPLICATE"
														case strings.Contains(Req.ResponseDetails.Body, "ALREADY_REGISTERED"):
															status = "ALREADY_REGISTERED"
															//InvalidAccs(Config.Email + ":" + Config.Password)
														case strings.Contains(Req.ResponseDetails.Body, "NOT_ENTITLED"):
															status = "NOT_ENTITLED"
															//InvalidAccs(Config.Email + ":" + Config.Password)
														default:
															switch Req.ResponseDetails.StatusCode {
															case "429":
																status = "RATE_LIMITED"
																//proxy = utils.Proxy.CompRand()
															case "200":
																if utils.Con.SkinChange.Link != "" {
																	apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Config.Bearer)
																}
																if utils.Con.Bools.UseWebhook {
																	go func() {
																		json, _ := BuildWebhook(name, "0", "")
																		err, ok := webhook.Webhook(utils.Con.WebhookURL, json)
																		if err != nil {
																			fmt.Println(utils.Logo(err.Error()))
																		} else if ok {
																			fmt.Println(utils.Logo("Succesfully sent personal webhook!"))
																		}
																	}()
																}
																NMC := utils.Namemc_key(Config.Bearer)
																if utils.Con.NMC.UseNMC {
																	Claim_NAMEMC(NMC)
																	SendFollow(name)
																}
																EmailClaimed = utils.Logo(fmt.Sprintf("✓ %v claimed %v @ %v [%v]\n", Config.Email, name, time.Now().Format("05.0000"), NMC))
																cl = true
																status = "CLAIMED"
															}
														}
														fmt.Println(utils.Logo(fmt.Sprintf(`✗ <%v> [%v] %v %v ⑇ %v ↪ %v`, time.Now().Format("15:04:05.0000"), Req.ResponseDetails.StatusCode, status, name, utils.HashEmailClean(Config.Email), strings.Split(c.Proxy, ":")[0])))
													}()
												}
												wg.Wait()
											}
										}(accs, c)
									}
								}
								Spread := time.Millisecond
								if utils.Con.Settings.UseCustomSpread {
									Spread = time.Duration(utils.Con.Settings.Spread) * time.Millisecond
								} else {
									Spread = TempCalc(utils.Con.Settings.SleepAmtPerGc, utils.Accamt)
								}
								time.Sleep(Spread)
							}
						}
					}()

					for {
						if cl || Changed {
							if EmailClaimed == "" {
								EmailClaimed = utils.Logo("No account has sniped the name.")
							}
							fmt.Printf(EmailClaimed)
							break
						}
						time.Sleep(1 * time.Second)
					}
				},
				Args: []string{
					"-u",
					"-start",
					"-end",
				},
			},
			"find_name": {
				Action: func() {
					if len(utils.Con.Bearers) != 0 {
						var detected bool
						var detected_string string
						var wg sync.WaitGroup
						var Name string = StrCmd.String("-name")
						body, _ := os.ReadFile("accounts.txt")
						scanner := bufio.NewScanner(bytes.NewBuffer(body))
						var accs []string
						for scanner.Scan() {
							accs = append(accs, scanner.Text())
						}

						var use_proxy, gcamt int
						var gc_a int
						var fgc bool = true
						var data map[string][]utils.Proxys_Accs = make(map[string][]utils.Proxys_Accs)
						for _, bearer := range accs {
							P := strings.Split(bearer, ":")
							if use_proxy >= len(utils.Proxy.Proxys) && len(utils.Proxy.Proxys) < len(accs) {
								break
							}
							if fgc {
								data["Account"] = []utils.Proxys_Accs{{Proxy: utils.Proxy.Proxys[use_proxy]}}
								fgc = false
								use_proxy++
							}
							var am int = 3
							if len(data["Account"][gc_a].Accs) != am {
								data["Account"][gc_a].Accs = append(data["Account"][gc_a].Accs, apiGO.Info{
									Email: P[0], Password: P[1],
								})
								gcamt++
							} else {
								gc_a++
								gcamt++
								data["Account"] = append(data["Account"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{{
									Email: P[0], Password: P[1],
								}}})
								use_proxy++
							}
						}
						for _, acc := range data["Account"] {
							ip, port, user, pass := "", "", "", ""
							switch data := strings.Split(acc.Proxy, ":"); len(data) {
							case 2:
								ip = data[0]
								port = data[1]
							case 4:
								ip = data[0]
								port = data[1]
								user = data[2]
								pass = data[3]
							}
							for _, acc := range acc.Accs {
								wg.Add(1)
								go func(acc apiGO.Info, ip, port, user, pass string) {
									defer wg.Done()
									info := apiGO.MS_authentication(acc.Email, acc.Password, &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass})
									if info.Info.Name == "" {
										info.Info.Name = "null"
									}
									if strings.EqualFold(info.Info.Name, Name) {
										detected = true
										detected_string = utils.Logo(fmt.Sprintf("<%v> %v - %v >> %v, %v", time.Now().Format("05.000"), info.Email, info.Password, info.Info.Name, info.AccountType))
									}
									fmt.Println(utils.Logo(fmt.Sprintf("<%v> %v >> %v, %v", time.Now().Format("05.000"), utils.HashEmailClean(info.Email), info.Info.Name, info.AccountType)))
								}(acc, ip, port, user, pass)
							}
						}
						wg.Wait()
						if detected {
							fmt.Println(detected_string)
						}
					}
				},
				Args: []string{
					"-name",
				},
			},
		},
	}
	if len(os.Args) > 1 {
		app.Input(strings.Join(os.Args[1:], " "))
	} else {
		app.Run()
	}
}

type Images struct {
	Image image.Image
	Url   string
	Row   int
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

func ReturnPayload(acc, bearer, name string) string {
	if acc == "Giftcard" {
		var JSON string = fmt.Sprintf(`{"profileName":"%v"}`, name)
		if utils.Con.Bools.UseMethod {
			return fmt.Sprintf("POST https://minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net/minecraft/profile// HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		} else {
			return fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nContent-Length:%v\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %v\r\n\r\n%v\r\n", len(JSON), bearer, JSON)
		}
	} else {
		return "PUT /minecraft/profile/name/" + name + " HTTP/1.1\r\nHost: minecraftapi-bef7bxczg0amd8ef.z01.azurefd.net\r\nConnection: open\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: Bearer " + bearer + "\r\n"
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
					if page, err := browser.NewPage(playwright.BrowserNewContextOptions{
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
									fmt.Println(utils.Logo(fmt.Sprintf("<%v> %v page required email verification to login..", time.Now().Format("05.000"), email)))
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

															fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Added recovery onto "+email+" >> "+email_code, time.Now().Format("05.000"))))
															utils.Con.Recovery = append(utils.Con.Recovery, utils.Succesful{
																Email:     email,
																Recovery:  email_code,
																Code_Used: code,
															})
															utils.Con.SaveConfig()
															utils.Con.LoadState()
															page.Close()
															break Exit
														}
													}
												}
												fmt.Println(utils.Logo(fmt.Sprintf("<%v> Attempt %v/100 to find a email code for %v >> %v", time.Now().Format("05.000"), i, email, email_code)))
												time.Sleep(5 * time.Second)
											}
										}
									}
								} else {
									if utils.Con.Bools.ApplyNewRecoveryToExistingEmails {
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
														fmt.Println(utils.Logo(fmt.Sprintf("<%v> i Changed recovery for "+email+" >> "+email_code, time.Now().Format("05.000"))))
														utils.Con.Recovery = append(utils.Con.Recovery, utils.Succesful{
															Email:     email,
															Recovery:  email_code,
															Code_Used: "nil",
														})
														utils.Con.SaveConfig()
														utils.Con.LoadState()
														page.Close()
														break Exit2
													}
												}
											}
											fmt.Println(utils.Logo(fmt.Sprintf("<%v> Attempt %v/100 to find a email code for %v >> %v", time.Now().Format("05.000"), i, email, email_code)))
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
			fmt.Println(utils.Logo(err.Error()))
		}
	} else {
		fmt.Println(utils.Logo(err.Error()))
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

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

func X_Forwarded_For() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%v.%v.%v.%v", rand.Intn(300-256)+256, rand.Intn(255), rand.Intn(255), rand.Intn(9))
}

func BuildWebhook(name, searches, headurl string) ([]byte, webhook.Web) {
	new := utils.Con.Webhook
	for i := range new.Embeds {
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{name}", name, -1)
		new.Embeds[i].Description = strings.Replace(new.Embeds[i].Description, "{searches}", searches, -1)
		new.Embeds[i].Author.Name = strings.Replace(new.Embeds[i].Author.Name, "{name}", name, -1)
		new.Embeds[i].Author.Name = strings.Replace(new.Embeds[i].Author.Name, "{searches}", searches, -1)
		new.Embeds[i].Author.IconURL = strings.Replace(new.Embeds[i].Author.IconURL, "{headurl}", headurl, -1)
		new.Embeds[i].Author.IconURL = strings.Replace(new.Embeds[i].Author.IconURL, "{name}", name, -1)
		new.Embeds[i].Author.URL = strings.Replace(new.Embeds[i].Author.URL, "{headurl}", headurl, -1)
		new.Embeds[i].Author.URL = strings.Replace(new.Embeds[i].Author.URL, "{name}", name, -1)
		new.Embeds[i].URL = strings.Replace(new.Embeds[i].URL, "{name}", name, -1)
		new.Embeds[i].Footer.Text = strings.Replace(new.Embeds[i].Footer.Text, "{name}", name, -1)
		new.Embeds[i].Footer.Text = strings.Replace(new.Embeds[i].Footer.Text, "{searches}", searches, -1)
		new.Embeds[i].Footer.IconURL = strings.Replace(new.Embeds[i].Footer.IconURL, "{name}", name, -1)
		new.Embeds[i].Footer.IconURL = strings.Replace(new.Embeds[i].Footer.IconURL, "{headurl}", headurl, -1)
		for e, field := range new.Embeds[i].Fields {
			field.Name = strings.Replace(field.Name, "{headurl}", headurl, -1)
			field.Name = strings.Replace(field.Name, "{searches}", searches, -1)
			field.Name = strings.Replace(field.Name, "{name}", name, -1)
			field.Value = strings.Replace(field.Value, "{headurl}", headurl, -1)
			field.Value = strings.Replace(field.Value, "{searches}", searches, -1)
			field.Value = strings.Replace(field.Value, "{name}", name, -1)
			new.Embeds[i].Fields[e] = field
		}
	}
	json, _ := json.Marshal(new)
	return json, new
}

func ReturnEmbed(name, searches, headurl string) (Data discordgo.MessageSend) {
	_, new := BuildWebhook(name, searches, headurl)
	for _, com := range new.Embeds {
		var Footer discordgo.MessageEmbedFooter
		var Image discordgo.MessageEmbedImage
		var Thumbnail discordgo.MessageEmbedThumbnail
		var Author discordgo.MessageEmbedAuthor

		if !reflect.DeepEqual(com.Footer, webhook.Footer{}) {
			Footer = discordgo.MessageEmbedFooter{
				Text:    com.Footer.Text,
				IconURL: com.Footer.IconURL,
			}
		}
		if !reflect.DeepEqual(com.Image, webhook.Image{}) {
			Image = discordgo.MessageEmbedImage{
				URL: com.Image.URL,
			}
		}
		if !reflect.DeepEqual(com.Thumbnail, webhook.Thumbnail{}) {
			Thumbnail = discordgo.MessageEmbedThumbnail{
				URL: com.Thumbnail.URL,
			}
		}
		if !reflect.DeepEqual(com.Author, webhook.Author{}) {
			Author = discordgo.MessageEmbedAuthor{
				URL:     com.Author.URL,
				Name:    com.Author.Name,
				IconURL: com.Author.IconURL,
			}
		}

		Data.Embeds = append(Data.Embeds, &discordgo.MessageEmbed{
			URL:         com.URL,
			Description: com.Description,
			Color:       com.Color,
			Footer:      &Footer,
			Image:       &Image,
			Thumbnail:   &Thumbnail,
			Author:      &Author,
			Fields:      returnjustfields(com),
		})
	}
	return
}

func returnjustfields(com webhook.Embeds) (Data []*discordgo.MessageEmbedField) {
	for _, c := range com.Fields {
		Data = append(Data, &discordgo.MessageEmbedField{
			Name:   c.Name,
			Value:  c.Value,
			Inline: c.Inline,
		})
	}
	return
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

func getProfiles() {
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
