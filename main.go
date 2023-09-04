package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"main/packages/utils"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/packages/Center"
	"main/packages/StrCmd"
	"main/packages/apiGO"

	"main/packages/utils/followbot"

	html "github.com/antchfx/htmlquery"
	"github.com/dop251/goja"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/term"
)

var (
	App_ = Center.App(`rgb(125,110,221)`, `rgb(90%,45%,97%)`, `hsl(229,79%,85%)`)
	logo = `,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,.....,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,:;;;;;;::,..,,,,,,,,,,,,
,,,,,,,,,,,:..1ftfffffffttti.,,,,,,,,,,,
,,,,,,,,,,,,::,itfLLfttfCCtCf.,,,,,,,,,,
,,,,,,,,,,,,,:,,LLLG,::1fCGiGi.,,,,,,,,,
,,,,,,,,,,,,.,,,i1LG:t :1t8it1.,,,,,,,,,
,,,,,,,,,,,,.,,.1LLG:1..if0;t;.,,,,,,,,,
,,,,,,,,,,,,,,. ;fLL.,.:1C1;i.,,,,,,,,,,
,,,,,,,,,,,,.,1itfLL11i11;;;,,,,,,,,,,,,
,,,,,,,,,,,,.,it11iii;:,,,,.,,,,,,,,,,,,
,,,,,,,,,,,,,,........,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
;,.,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,`
)

type Images struct {
	Image image.Image
	Url   string
	Row   int
}

func init() {
	if _, err := os.Stat("/data/yt"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("data/yt", os.ModePerm)
	}

	if files, err := os.ReadDir("data/yt"); err == nil {
		for _, file := range files {
			os.Remove("data/yt/" + file.Name())
			// please clear ur recycle bin when needed.
		}
	}

	apiGO.Clear()
	utils.Con.LoadState()

	go YTPlayer(utils.Con.Settings.Youtube)

	utils.Roots.AppendCertsFromPEM(utils.ProxyByte)

	Logo()

	// get screen size
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	var values string
	scanners := bufio.NewScanner(strings.NewReader(`
┏┓       ┓ ┓
┃ ┏┓┓┏┏┳┓┣┓┃┏┓
┗┛┛ ┗┻┛┗┗┗┛┗┗`))
	for scanners.Scan() {
		w_ := scanners.Text()
		values += strings.Repeat(" ", w/2-6) + w_ + "\r\n"
	}

	if len(utils.RGB) != 0 {
		App_.Grad.Mutline(values)
	}
	fmt.Print(App_.Grad.Mutline(values))

	var AccountsVer []string
	file, _ := os.Open("data/accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	// Pre check accounts.txt
	if len(AccountsVer) == 0 {
		fmt.Println(utils.Logo("\n" + `[ERROR] Unable to continue, no accounts inside of data/accounts.txt` + "\n"))
		return
	}

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

	if utils.Con.Bools.FirstUse {
		fmt.Print(utils.Logo("Use proxys for authentication? : [YES/NO] > "))
		var ProxyAuth string
		fmt.Scan(&ProxyAuth)
		utils.Con.Bools.FirstUse = false
		utils.Con.Bools.UseProxyDuringAuth = strings.Contains(strings.ToLower(ProxyAuth), "y")
		utils.Con.SaveConfig()
		utils.Con.LoadState()
	}

	if _, err := os.Stat("/data/skinarts"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("data/skinarts", os.ModePerm)
	}
	if _, err := os.Stat("/data/frontpagedata"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("data/frontpagedata", os.ModePerm)
	}

	if _, err := os.Stat("data"); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll("data", os.ModePerm)
		os.Create("data/config.json")
		os.Create("data/accounts.txt")
		os.Create("data/proxys.txt")
		os.Create("data/invalids.txt")
	}

	utils.Proxy.GetProxys(false, nil)
	utils.Proxy.Setup()
	utils.AuthAccs()

	if len(utils.Bearer.Details) == 0 || len(utils.Con.Bearers) == 0 {
		fmt.Println(utils.Logo("\n" + "[ERROR] No bearers appeared to have authed, refer to data/invalids.txt" + "\n"))
	}

	utils.Regenerateallaccs()
	go utils.CheckAccs()
	var use_proxy, gcamt, mfaamt int

	for _, bearer := range utils.Bearer.Details {
		if use_proxy >= len(utils.Proxy.Proxys) && len(utils.Proxy.Proxys) < len(utils.Bearer.Details) {
			break
		}
		switch bearer.AccountType {
		case "Microsoft":
			if utils.First_mfa {
				utils.Accs["Microsoft"] = []utils.Proxys_Accs{{Proxy: utils.Proxy.Proxys[use_proxy]}}
				utils.First_mfa = false
				use_proxy++
			}
			var am int = utils.Con.Settings.AccountsPerMfa
			if am == 0 {
				am = 1
			}
			if len(utils.Accs["Microsoft"][utils.Use_mfa].Accs) != am {
				utils.Accs["Microsoft"][utils.Use_mfa].Accs = append(utils.Accs["Microsoft"][utils.Use_mfa].Accs, bearer)
				utils.Accamt++
				mfaamt++
			} else {
				utils.Use_mfa++
				utils.Accamt++
				mfaamt++
				utils.Accs["Microsoft"] = append(utils.Accs["Microsoft"], utils.Proxys_Accs{Proxy: utils.Proxy.Proxys[use_proxy], Accs: []apiGO.Info{bearer}})
				use_proxy++
			}
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
			followbot.Cookie = utils.Con.CF.Tokens
			go func() {
				for {
					time.Sleep(time.Until(time.Unix(utils.Con.CF.GennedAT, 0)))
					if status, cookies := followbot.Get_CF_Clearance(); status == 200 {
						for _, cookies := range cookies {
							if cookies.Name == "cf_clearance" {
								utils.Con.CF = utils.CF{
									Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
									GennedAT: time.Now().Add(time.Second * 1800).Unix(),
								}
								followbot.Cookie = utils.Con.CF.Tokens
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
			if status, cookies := followbot.Get_CF_Clearance(); status == 200 {
				for _, cookies := range cookies {
					if cookies.Name == "cf_clearance" {
						utils.Con.CF = utils.CF{
							Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
							GennedAT: time.Now().Add(time.Second * 1800).Unix(),
						}
						followbot.Cookie = utils.Con.CF.Tokens
						utils.Con.SaveConfig()
						utils.Con.LoadState()
						break
					}
				}
				go func() {
					for {
						time.Sleep(time.Until(time.Unix(utils.Con.CF.GennedAT, 0)))
						if status, cookies := followbot.Get_CF_Clearance(); status == 200 {
							for _, cookies := range cookies {
								if cookies.Name == "cf_clearance" {
									utils.Con.CF = utils.CF{
										Tokens:   fmt.Sprintf("%v=%v", cookies.Name, cookies.Value),
										GennedAT: time.Now().Add(time.Second * 1800).Unix(),
									}
									followbot.Cookie = utils.Con.CF.Tokens
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
				req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
				req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Add("Origin", "https://namemc.com")
				req.Header.Add("Referer", "https://namemc.com/login")
				Cl := http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						for _, name := range req.Response.Cookies() {
							if name.Name == "session_id" {
								followbot.Session.Key = name.Value
							}
						}
						return nil
					},
				}
				r, _ := Cl.Do(req)
				if r.StatusCode == 401 {
					fmt.Println(utils.Logo(fmt.Sprintf("<%v> Email and password invalid for namemc..", time.Now().Format("05.000"))))
				} else {
					if followbot.Session.Key != "" {
						utils.Con.NMC.NamemcLoginData = utils.NMC{
							Token:      followbot.Session.Key,
							LastAuthed: time.Now().Add(time.Second * 86400).Unix(),
						}
						utils.Con.SaveConfig()
						utils.Con.LoadState()
					}
				}
			}
		} else {
			followbot.Session.Key = utils.Con.NMC.NamemcLoginData.Token
		}
		followbot.GetProfiles()
	}
}

func main() {
	app := StrCmd.App{
		Display:        utils.Logo("⋅ "),
		Version:        "v1.16.1-BETA",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"namemc": {
				Description: "The command that can handle FUN commands related to NAMEMC!",
				Subcommand: map[string]StrCmd.SubCmd{
					"info": {
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							req, _ := http.NewRequest("GET", "https://namemc.com/search?q="+StrCmd.String("-name"), nil)
							req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
							req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
							req.Header.Add("Origin", "https://namemc.com")
							req.Header.Add("X-Forwarded-For", followbot.X_Forwarded_For())
							req.Header.Add("Referer", "https://namemc.com/")
							if resp, err := http.DefaultClient.Do(req); err == nil {
								resp_body, _ := io.ReadAll(resp.Body)
								namemc := followbot.GetInfo(resp_body)
								fmt.Println(utils.Logo(fmt.Sprintf(`
   Name: %v
  Views: %v
HeadURL: %v
BodyURL: %v
  Start: %v
    End: %v
 Status: %v
`, StrCmd.String("-name"), namemc.Searches, namemc.HeadURL, namemc.BodyUrl, namemc.Start, namemc.End, namemc.Status)))
							}
						},
						Args: []string{"-name"},
					},
					"get-skins": {
						Args: []string{"-pages"},
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							amt := StrCmd.Int("-pages")
							var Path = ""
							var Skins []followbot.NamemcSkins
							for i := 2; i < amt+2; i++ {
								req, _ := http.NewRequest("GET", "https://namemc.com/minecraft-skins"+Path, nil)
								req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
								req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
								req.Header.Add("Origin", "https://namemc.com")
								req.Header.Add("X-Forwarded-For", followbot.X_Forwarded_For())
								req.Header.Add("Referer", "https://namemc.com/minecraft-skins")
								resp, _ := http.DefaultClient.Do(req)
								if resp.StatusCode == 200 {
									resp_body, _ := io.ReadAll(resp.Body)
									Skins = append(Skins, followbot.GetAllSkinsOnPage(string(resp_body))...)
									Path = fmt.Sprintf("?page=%v", i)
								} else {
									fmt.Println(resp.StatusCode)
								}
							}
							path := fmt.Sprintf("data/frontpagedata/%v", time.Now().Unix())
							os.Mkdir(path, 0755)
							for _, page := range Skins {
								if resp, err := http.Get(page.SkinDownload); err == nil {
									if img, err := png.Decode(resp.Body); err == nil {
										file, _ := os.Create(path + "/" + page.Number + ".png")
										png.Encode(file, img)
										file.Close()
									} else {
										fmt.Println(err)
									}
								} else {
									fmt.Println(err)
								}
							}
						},
					},
					"follow": {
						Args: []string{"-name"},
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							Name := StrCmd.String("-name")
							req, _ := http.NewRequest("GET", "https://namemc.com/profile/"+Name, nil)
							req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
							req.AddCookie(&http.Cookie{Name: "session_id", Value: followbot.Session.Key})
							req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
							var redirect string
							Client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { redirect = req.URL.String(); return nil }}
							resp, _ := Client.Do(req)
							if resp.StatusCode == 200 {
								body, _ := io.ReadAll(resp.Body)
								var Target followbot.Target = followbot.Target{
									UUID: strings.Split(strings.Split(string(body), `order-lg-2 col-lg" style="font-size: 90%"><samp>`)[1], `</samp></div>`)[0],
									URL:  redirect,
								}
								for _, acc := range followbot.Session.Accounts {
									Target.Send()
									if acc.Name != followbot.Session.DisplayName {
										followbot.Swapprofile(acc.URLPath)
									}
								}
							}
						},
					},
					"rl-test": {
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							type Resps struct {
								Name       string     `json:"name,omitempty" bson:"name"`
								Start      *time.Time `json:"start_date,omitempty"`
								End        *time.Time `json:"end_date,omitempty"`
								Start_Unix int64      `json:"start_unix,omitempty"`
								End_Unix   int64      `json:"end_unix,omitempty"`
								CachedAt   string     `json:"cachedat,omitempty" bson:"cachedat"`
								Searches   string     `json:"searches,omitempty" bson:"searches"`
							}
							var Info []Resps
							next := "https://namemc.com/minecraft-names"
						Exit:
							for {
								var I []Resps

								req, _ := http.NewRequest("GET", next, nil)
								req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
								req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

								IP := followbot.X_Forwarded_For()
								req.Header.Add("X-Forwarded-For", IP)

								if resp, err := http.DefaultClient.Do(req); err == nil {
									if resp.StatusCode == 200 {
										data, _ := io.ReadAll(resp.Body)
										if doc, err := html.Parse(strings.NewReader(string(data))); err == nil {
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
												start := html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td[2]/time", i))
												end := html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td/span", i))
												if start != nil && end != nil {
													if len(start.Attr) > 0 && len(end.Attr) > 0 {
														t1, _ := time.Parse(time.RFC3339, start.Attr[0].Val)
														t2, _ := time.Parse(time.RFC3339, end.Attr[2].Val)
														I = append(I, Resps{
															Name:       strings.ToLower(html.FindOne(doc, fmt.Sprintf("/html/body/main/div/div[4]/div/table/tbody/tr[%v]/td[1]/a", i)).FirstChild.Data),
															Start:      &t1,
															Start_Unix: t1.Unix(),
															End:        &t2,
															End_Unix:   t2.Unix(),
															Searches:   Searches,
														})
													}
												}
											}
											Info = append(Info, I...)
											fmt.Printf("<%v> %v collected %v name info...\n", time.Now().Format("05.000"), IP, len(Info))
											next = fmt.Sprintf("https://namemc.com%v", html.FindOne(doc, "/html/body/main/div/div[5]/nav/ul/li[4]/a").Attr[1].Val)
											if len(I) < 3 {
												break Exit
											}
										}
									}
								}
							}
							j, _ := json.MarshalIndent(Info, " ", "  ")
							os.Create("namemc.json")
							os.WriteFile("namemc.json", j, 0644)
						},
					},
					"skinart-png": {
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							var name string
							fmt.Print(utils.Logo("Name of the profile you wanna scrape: "))
							fmt.Scan(&name)
							resp, _ := http.Get("https://namemc.info/data/namemc/skinart/logo/" + name)
							img, _, _ := image.Decode(resp.Body)
							path := "data/skinarts/" + strings.ReplaceAll(uuid.NewString(), "-", "") + "_" + name + ".png"
							out, _ := os.Create(path)
							png.Encode(out, img)
						},
					},
				},
			},
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
			"recover": {Action: func() {
				force := StrCmd.Bool("--force")
				use := []apiGO.Info{}
				if force {
					use = utils.Bearer.Details
				} else {
					if _, err := os.Stat("data/invalids.txt"); !os.IsNotExist(err) {
						body, _ := os.ReadFile("data/invalids.txt")
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
				utils.AddRecoveryEmails(use)
			}, Args: []string{"--force"}},
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					if len(utils.Bearer.Details) != 0 {
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
															Req := apiGO.Details{ResponseDetails: apiGO.SocketSending(P, utils.ReturnPayload(Config.AccountType, Config.Bearer, name)), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType}
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

																	NMC := utils.Namemc_key(Config.Bearer)
																	if utils.Con.NMC.UseNMC {
																		followbot.Claim_NAMEMC(NMC)
																		followbot.SendFollow(name)
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
					} else {
						if len(utils.Con.Bearers) == 0 && len(utils.Bearer.Details) == 0 {
							return
						}
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
						body, _ := os.ReadFile("data/accounts.txt")
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
		if followbot.Session.Key != "" {
			app.Display = utils.Logo(followbot.Session.DisplayName + " ⋅ ")
		}
		app.Run()
	}
}

func TempCalc(interval, accamt int) time.Duration {
	amt := interval / accamt
	if amt < 10 {
		return time.Duration(amt) * time.Millisecond
	}
	return time.Duration(amt) * time.Millisecond
}

func Logo() {

	App_.PrintMiddleUncachedToBody(logo)

	time.Sleep(5 * time.Second)

	Center.Clear()

}

type YT struct {
	StreamingData StreamingData `json:"streamingData"`
}
type Formats struct {
	Itag             int    `json:"itag"`
	URL              string `json:"url"`
	MimeType         string `json:"mimeType"`
	Bitrate          int    `json:"bitrate"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	LastModified     string `json:"lastModified"`
	Quality          string `json:"quality"`
	Fps              int    `json:"fps"`
	QualityLabel     string `json:"qualityLabel"`
	ProjectionType   string `json:"projectionType"`
	AudioQuality     string `json:"audioQuality"`
	ApproxDurationMs string `json:"approxDurationMs"`
	AudioSampleRate  string `json:"audioSampleRate"`
	AudioChannels    int    `json:"audioChannels"`
	Sig              string `json:"signatureCipher"`
}
type InitRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type IndexRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type ColorInfo struct {
	Primaries               string `json:"primaries"`
	TransferCharacteristics string `json:"transferCharacteristics"`
	MatrixCoefficients      string `json:"matrixCoefficients"`
}
type AdaptiveFormats struct {
	Itag             int        `json:"itag"`
	URL              string     `json:"url"`
	MimeType         string     `json:"mimeType"`
	Bitrate          int        `json:"bitrate"`
	Width            int        `json:"width,omitempty"`
	Height           int        `json:"height,omitempty"`
	InitRange        InitRange  `json:"initRange"`
	IndexRange       IndexRange `json:"indexRange"`
	LastModified     string     `json:"lastModified"`
	ContentLength    string     `json:"contentLength"`
	Quality          string     `json:"quality"`
	Fps              int        `json:"fps,omitempty"`
	QualityLabel     string     `json:"qualityLabel,omitempty"`
	ProjectionType   string     `json:"projectionType"`
	AverageBitrate   int        `json:"averageBitrate"`
	ApproxDurationMs string     `json:"approxDurationMs"`
	ColorInfo        ColorInfo  `json:"colorInfo,omitempty"`
	HighReplication  bool       `json:"highReplication,omitempty"`
	AudioQuality     string     `json:"audioQuality,omitempty"`
	AudioSampleRate  string     `json:"audioSampleRate,omitempty"`
	AudioChannels    int        `json:"audioChannels,omitempty"`
	LoudnessDb       float64    `json:"loudnessDb,omitempty"`
	Sig              string     `json:"signatureCipher"`
}
type StreamingData struct {
	ExpiresInSeconds string            `json:"expiresInSeconds"`
	Formats          []Formats         `json:"formats"`
	AdaptiveFormats  []AdaptiveFormats `json:"adaptiveFormats"`
}

var Y = regexp.MustCompile(`^.*(youtu.be\/|v\/|e\/|u\/\w+\/|embed\/|v=)([^#\&\?]*).*`)

func YoutubeURL(URL string) string {
	YT_ := Y.FindAllStringSubmatch(URL, -1)
	if len(YT_) > 0 {
		if len(YT_[0]) > 2 {
			return YT_[0][2]
		}
	}
	return "Unknown"
}

func YTPlayer(url string) {
	if url == "" {
		return
	}
	var IDs []string
	if strings.Contains(url, "playlist?list=") {
		req, _ := http.NewRequest("GET", url, nil)
		rr, _ := http.DefaultClient.Do(req)
		aa, _ := io.ReadAll(rr.Body)
		var DD utils.YTPageConfig
		json.Unmarshal([]byte(strings.Split(strings.Split(string(aa), `var ytInitialData =`)[1], `;</script>`)[0]), &DD)
		for _, data := range DD.Contents.TwoColumnBrowseResultsRenderer.Tabs {
			for _, yt := range data.TabRenderer.Content.SectionListRenderer.Contents {
				for _, pagedata := range yt.ItemSectionRenderer.Contents {
					for _, data := range pagedata.PlaylistVideoListRenderer.Contents {
						IDs = append(IDs, data.PlaylistVideoRenderer.VideoID)
					}
				}
			}
		}
	} else {
		ID := YoutubeURL(url)
		IDs = append(IDs, ID)
	}
	GetSongs(IDs)
}

type YTVids struct {
	Body  []byte
	Index int
}

func GetSongs(IDs []string) {

	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		if runtime.GOOS == "windows" {
			if err := exec.Command("winget", "install", "ffmpeg").Run(); err != nil {
				fmt.Println(utils.Logo(err.Error() + " - Closing songs(s).."))
				return
			} else {
				ffmpeg, _ = exec.LookPath("ffmpeg")
			}
		}
		return
	}

	for _, ID := range IDs {
		Body := fmt.Sprintf(`
		{
		  "context": {
			"client": {
			  "clientName": "WEB",
			  "clientVersion": "2.20230615.02.01"
			}
		  },
		  "videoId": "%v"
		}
		`, ID)
		req, _ := http.NewRequest("POST", "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false", bytes.NewReader([]byte(Body)))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(Body)))
		resp, e := http.DefaultClient.Do(req)
		if e == nil {
			defer resp.Body.Close()
			var Test YT
			get_id, _ := io.ReadAll(resp.Body)
			if strings.Contains(string(get_id), "This is a private video. Please sign in to verify that you may see it.") {
				return
			} else {
				json.Unmarshal(get_id, &Test)
				var found bool
				var URL, S, SIG string
				var ContentL int
				var sigtriggered bool
				for _, T := range Test.StreamingData.AdaptiveFormats {
					if T.AudioQuality == "AUDIO_QUALITY_LOW" {
						if T.URL != "" {
							found = true
							URL = T.URL
							if int_value, err := strconv.Atoi(T.ContentLength); err == nil {
								ContentL = int_value
							}
							break
						}
						if T.Sig != "" {
							URL = strings.Split(T.Sig, "url=")[1]
							SIG = T.Sig
							S = strings.Split(strings.Split(T.Sig, "s=")[1], "&sp=")[0]
							if int_value, err := strconv.Atoi(T.ContentLength); err == nil {
								ContentL = int_value
							}
							sigtriggered = true
							break
						}
					}
				}

				if !found {
					for _, T := range Test.StreamingData.Formats {
						if T.AudioQuality == "AUDIO_QUALITY_LOW" && T.URL != "" {
							found = true
							URL = T.URL
							resp, _ := http.Get(T.URL)
							if int_value, err := strconv.Atoi(resp.Header.Get("Content-Length")); err == nil {
								ContentL = int_value
							}
							break
						}
					}
				}

				if sigtriggered && URL != "" {

					pars, _ := url.ParseQuery(SIG)

					u, _ := url.Parse(pars.Get("url"))

					S, _ = url.PathUnescape(S)
					a, _ := decrypt([]byte(S), ID)
					S = string(a)
					q := u.Query()

					bb, _ := getPlayerConfig(ID)

					q.Add(pars.Get("sp"), S)

					vals, _ := decryptNParam(bb, q)
					u.RawQuery = vals.Encode()
					URL = u.String()

					resp, _ := http.Get(URL)
					ContentL, _ = strconv.Atoi(resp.Header.Get("Content-Length"))

					found = true
				}

				if found {
					var (
						seperated_values = ContentL / 2
						start_pos        = 0
						Data             = []YTVids{}
						wg               = sync.WaitGroup{}
						inp              = "data/yt/" + ID + ".mp3"
						out              = fmt.Sprintf("data/yt/audio_%v_out.mp3", time.Now().UnixNano())
					)

					file, _ := os.Create(inp)

					for i := 0; i < 2; i++ {
						wg.Add(1)
						if i == 2 {
							seperated_values = ContentL
						}
						go func(start, end, index int) {
							if sigtriggered {
								resp, _ := http.Get(URL + fmt.Sprintf("&range=%v-%v", start, end))
								Data = append(Data, YTVids{
									Body:  []byte(apiGO.ReturnJustString(io.ReadAll(resp.Body))),
									Index: index,
								})
							} else {
								req, _ := http.NewRequest("GET", URL, nil)
								req.Header.Add("Accept", "*/*")
								req.Header.Add("Range", fmt.Sprintf(fmt.Sprintf("bytes=%v-%v", start, end)))
								req.Header.Add("Referer", URL)
								r, _ := http.DefaultClient.Do(req)
								Data = append(Data, YTVids{
									Body:  []byte(apiGO.ReturnJustString(io.ReadAll(r.Body))),
									Index: index,
								})
							}
							wg.Done()
						}(start_pos, seperated_values, i)
						start_pos = seperated_values
						seperated_values = seperated_values + ContentL/2
					}

					wg.Wait()

					sort.Slice(Data, func(i, j int) bool {
						return Data[i].Index < Data[j].Index
					})

					var Vid []byte

					for _, body := range Data {
						Vid = append(Vid, body.Body...)
					}

					file.Write(Vid)
					file.Close()

					cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", inp, "-vn", out)
					cmd.Run()
					os.Remove(inp)
					file, err := os.Open(out)
					if err != nil {
						break
					}
					streamer, format, err := mp3.Decode(file)
					if err != nil {
						panic(err)
					}

					speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
					done := make(chan bool)
					speaker.Play(beep.Seq(streamer, beep.Callback(func() {
						done <- true
					})))
					<-done
					speaker.Clear()
					speaker.Close()
					file.Close()
					os.Remove(out)
				}
			}
		}
	}
}

type DecipherOperation func([]byte) []byte

func decrypt(cyphertext []byte, id string) ([]byte, error) {
	operations, err := parseDecipherOps(cyphertext, id)
	if err != nil {
		return nil, err
	}

	// apply operations
	bs := []byte(cyphertext)
	for _, op := range operations {
		bs = op(bs)
	}

	return bs, nil
}

const (
	jsvarStr   = "[a-zA-Z_\\$][a-zA-Z_0-9]*"
	reverseStr = ":function\\(a\\)\\{" +
		"(?:return )?a\\.reverse\\(\\)" +
		"\\}"
	spliceStr = ":function\\(a,b\\)\\{" +
		"a\\.splice\\(0,b\\)" +
		"\\}"
	swapStr = ":function\\(a,b\\)\\{" +
		"var c=a\\[0\\];a\\[0\\]=a\\[b(?:%a\\.length)?\\];a\\[b(?:%a\\.length)?\\]=c(?:;return a)?" +
		"\\}"
)

var (
	nFunctionNameRegexp = regexp.MustCompile("\\.get\\(\"n\"\\)\\)&&\\(b=([a-zA-Z0-9$]{0,3})\\[(\\d+)\\](.+)\\|\\|([a-zA-Z0-9]{0,3})")
	actionsObjRegexp    = regexp.MustCompile(fmt.Sprintf(
		"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s),?\\n?)+)\\};", jsvarStr, jsvarStr, swapStr, jsvarStr, spliceStr, jsvarStr, reverseStr))
	actionsFuncRegexp = regexp.MustCompile(fmt.Sprintf(
		"function(?: %s)?\\(a\\)\\{"+
			"a=a\\.split\\(\"\"\\);\\s*"+
			"((?:(?:a=)?%s\\.%s\\(a,\\d+\\);)+)"+
			"return a\\.join\\(\"\"\\)"+
			"\\}", jsvarStr, jsvarStr, jsvarStr))
	reverseRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, reverseStr))
	spliceRegexp  = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, spliceStr))
	swapRegexp    = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, swapStr))
)

func decryptNParam(config []byte, query url.Values) (url.Values, error) {
	// decrypt n-parameter
	nSig := query.Get("v")
	if nSig != "" {
		nDecoded, err := decodeNsig(config, nSig)
		if err != nil {
			return nil, fmt.Errorf("unable to decode nSig: %w", err)
		}
		query.Set("v", nDecoded)
	}

	return query, nil
}

func parseDecipherOps(config []byte, id string) (operations []DecipherOperation, err error) {
	config, _ = getPlayerConfig(id)
	objResult := actionsObjRegexp.FindSubmatch(config)
	funcResult := actionsFuncRegexp.FindSubmatch(config)
	if len(objResult) < 3 || len(funcResult) < 2 {
		return nil, fmt.Errorf("error parsing signature tokens (#obj=%d, #func=%d)", len(objResult), len(funcResult))
	}
	obj := objResult[1]
	objBody := objResult[2]
	funcBody := funcResult[1]
	var reverseKey, spliceKey, swapKey string
	if result := reverseRegexp.FindSubmatch(objBody); len(result) > 1 {
		reverseKey = string(result[1])
	}
	if result := spliceRegexp.FindSubmatch(objBody); len(result) > 1 {
		spliceKey = string(result[1])
	}
	if result := swapRegexp.FindSubmatch(objBody); len(result) > 1 {
		swapKey = string(result[1])
	}
	regex, err := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s|%s|%s)\\(a,(\\d+)\\)", regexp.QuoteMeta(string(obj)), regexp.QuoteMeta(reverseKey), regexp.QuoteMeta(spliceKey), regexp.QuoteMeta(swapKey)))
	if err != nil {
		return nil, err
	}
	var ops []DecipherOperation
	for _, s := range regex.FindAllSubmatch(funcBody, -1) {
		switch string(s[1]) {
		case reverseKey:
			ops = append(ops, reverseFunc)
		case swapKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, newSwapFunc(arg))
		case spliceKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, newSpliceFunc(arg))
		}
	}
	return ops, nil
}

var basejsPattern = regexp.MustCompile(`(/s/player/\w+/player_ias.vflset/\w+/base.js)`)

func getPlayerConfig(videoID string) ([]byte, error) {
	embedURL := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoID)
	req, _ := http.NewRequest("GET", embedURL, nil)
	req.Header.Set("Origin", "https://youtube.com")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	resp, _ := http.DefaultClient.Do(req)
	embedBody, _ := io.ReadAll(resp.Body)
	playerPath := string(basejsPattern.Find(embedBody))
	if playerPath == "" {
		return nil, errors.New("unable to find basejs URL in playerConfig")
	}
	reqa, _ := http.NewRequest("GET", "https://youtube.com"+playerPath, nil)
	reqa.Header.Set("Origin", "https://youtube.com")
	reqa.Header.Set("Sec-Fetch-Mode", "navigate")
	respa, _ := http.DefaultClient.Do(reqa)
	re, _ := io.ReadAll(respa.Body)
	return re, nil
}

func reverseFunc(bs []byte) []byte {
	l, r := 0, len(bs)-1
	for l < r {
		bs[l], bs[r] = bs[r], bs[l]
		l++
		r--
	}
	return bs
}

func newSwapFunc(arg int) DecipherOperation {
	return func(bs []byte) []byte {
		pos := arg % len(bs)
		bs[0], bs[pos] = bs[pos], bs[0]
		return bs
	}
}

func newSpliceFunc(pos int) DecipherOperation {
	return func(bs []byte) []byte {
		return bs[pos:]
	}
}

func getNFunction(config []byte) (string, error) {
	nameResult := nFunctionNameRegexp.FindSubmatch(config)
	if len(nameResult) == 0 {
		return "", errors.New("unable to extract n-function name")
	}

	var name string
	if idx, _ := strconv.Atoi(string(nameResult[2])); idx == 0 {
		name = string(nameResult[4])
	} else {
		name = string(nameResult[1])
	}

	return extraFunction(config, name)

}

func decodeNsig(config []byte, encoded string) (string, error) {
	fBody, err := getNFunction(config)
	if err != nil {
		return "", err
	}

	return evalJavascript(fBody, encoded)
}

func evalJavascript(jsFunction, arg string) (string, error) {
	const myName = "myFunction"

	vm := goja.New()
	_, err := vm.RunString(myName + "=" + jsFunction)
	if err != nil {
		return "", err
	}

	var output func(string) string
	err = vm.ExportTo(vm.Get(myName), &output)
	if err != nil {
		return "", err
	}

	return output(arg), nil
}

func extraFunction(config []byte, name string) (string, error) {
	// find the beginning of the function
	def := []byte(name + "=function(")
	start := bytes.Index(config, def)
	if start < 1 {
		return "", fmt.Errorf("unable to extract n-function body: looking for '%s'", def)
	}

	// start after the first curly bracket
	pos := start + bytes.IndexByte(config[start:], '{') + 1

	var strChar byte

	// find the bracket closing the function
	for brackets := 1; brackets > 0; pos++ {
		b := config[pos]
		switch b {
		case '{':
			if strChar == 0 {
				brackets++
			}
		case '}':
			if strChar == 0 {
				brackets--
			}
		case '`', '"', '\'':
			if config[pos-1] == '\\' && config[pos-2] != '\\' {
				continue
			}
			if strChar == 0 {
				strChar = b
			} else if strChar == b {
				strChar = 0
			}
		}
	}

	return string(config[start:pos]), nil
}
