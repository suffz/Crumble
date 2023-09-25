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
	"main/packages/h2"
	"main/packages/utils"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/packages/Center"
	"main/packages/StrCmd"
	"main/packages/apiGO"

	"main/packages/utils/followbot"

	html "github.com/antchfx/htmlquery"
	"github.com/faiface/beep/speaker"
	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"github.com/suffz/Youtube"
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

var CURRSONG string

func init() {

	if _, err := os.Stat("data"); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll("data", os.ModePerm)
		os.Create("data/config.json")
		os.Create("data/accounts.txt")
		os.Create("data/proxys.txt")
		os.Create("data/invalids.txt")
		os.Mkdir("data/skinarts", os.ModePerm)
		os.Mkdir("data/frontpagedata", os.ModePerm)
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
	for _, rgb := range utils.Con.Gradient {
		utils.RGB = append(utils.RGB, fmt.Sprintf("rgb(%v,%v,%v)", rgb.R, rgb.G, rgb.B))
	}

	utils.Roots.AppendCertsFromPEM(utils.ProxyByte)

	Logo()

	fmt.Print(utils.Logo(`/~' _   _ _ |_ | _ 
\_,||_|| | ||_)|(/_
`))

	var AccountsVer []string

	scanner := bufio.NewScanner(bytes.NewBufferString(apiGO.ReturnJustString(os.ReadFile("data/accounts.txt"))))
	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	// Pre check accounts.txt
	if len(AccountsVer) == 0 {
		fmt.Print(utils.Logo("\n" + `[ERROR] Unable to continue, no accounts inside of data/accounts.txt` + "\n"))
		return
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
				Headers := followbot.CommonHeaders
				Headers["origin"] = "https://namemc.com"
				Headers["referer"] = "https://namemc.com/login"
				Headers["content-type"] = "application/x-www-form-urlencoded"
				request, _ := h2.BuildRequest("https://namemc.com/login", "POST", fmt.Sprintf(`email=%v&password=%v`, acc[0], acc[1]), Headers, followbot.ReturnCF())
				if _, resp, err := request.Commit(nil); err == nil {
					switch resp.StatusCode {
					case 401:
						fmt.Println(utils.Logo(fmt.Sprintf("<%v> Email and password invalid for namemc..", time.Now().Format("05.000"))))
					default:
						for _, cookie := range resp.Cookies() {
							if cookie.Name == "session_id" {
								followbot.Session.Key = cookie.Value
								utils.Con.NMC.NamemcLoginData = utils.NMC{
									Token:      followbot.Session.Key,
									LastAuthed: time.Now().Add(time.Second * 86400).Unix(),
								}
								utils.Con.SaveConfig()
								utils.Con.LoadState()
							}
						}
					}
				}
			}
		} else {
			followbot.Session.Key = utils.Con.NMC.NamemcLoginData.Token
		}
		followbot.GetProfiles()
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

	utils.Proxy.GetProxys(false, nil)
	utils.Proxy.Setup()

	utils.AuthAccs()

	if len(utils.Bearer.Details) == 0 {
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
	if CURRSONG == "" && utils.Con.NMC.UseNMC && utils.Con.NMC.NamemcLoginData.Token != "" {
		CURRSONG = followbot.Session.DisplayName
	}

	if utils.Con.Settings.Youtube != "" {
		if strings.Contains(utils.Con.Settings.Youtube, "playlist?list=") {
			Reqs := Youtube.Playlist(utils.Con.Settings.Youtube)
			go func() {
				for {
					for _, req := range Reqs {
						CURRSONG = req.Title
						YT := Youtube.Video(req, false, Youtube.AudioMedium)
						if resp, _, err := YT.Download(); err == nil {
							YT.Play("data/yt/"+req.ID+".mp3", fmt.Sprintf("data/yt/audio_%v_out.mp3", time.Now().UnixNano()), resp)
						}
					}
				}
			}()
		} else {
			ID := Youtube.YoutubeURL(utils.Con.Settings.Youtube)
			YT := Youtube.Video(Youtube.Youtube{
				ID: ID,
			}, false, Youtube.AudioMedium)
			CURRSONG = YT.Config.Title
			go func() {
				if resp, _, err := YT.Download(); err == nil {
					s := YT.GetStream("data/yt/"+ID+".mp3", fmt.Sprintf("data/yt/audio_%v_out.mp3", time.Now().UnixNano()), resp)
					for {
						YT.PlayWithStream(s)
					}
				}
			}()
		}
	}
}

func main() {
	app := StrCmd.App{
		Display:        utils.Logo(CURRSONG + " ⋅ "),
		Version:        "v1.16.2",
		AppDescription: "Crumble is a open source minecraft turbo!",
		Commands: map[string]StrCmd.Command{
			"namemc": {
				Description: `The command that can handle FUN commands related to NAMEMC! (must have CF_CLEARANCE Token! "usenamemc_fordroptime_andautofollow": true in the config..)`,
				Subcommand: map[string]StrCmd.SubCmd{
					"info": {
						Description: "gets information on the -name you supply.",
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							namemc := followbot.Info(StrCmd.String("-name"))
							fmt.Println(utils.Logo(fmt.Sprintf(`
   Name: %v
  Views: %v
HeadURL: %v
BodyURL: %v
  Start: %v
    End: %v
 Status: %v
`, StrCmd.String("-name"), namemc.Searches, namemc.HeadURL, namemc.BodyUrl, namemc.Start, namemc.End, namemc.Status)))
						},
						Args: []string{"-name"},
					},
					"get-skins": {
						Description: "Takes popular skins downloads them and saves them to your frontpagedata folder.",
						Args:        []string{"-pages"},
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							amt := StrCmd.Int("-pages")
							var Path = ""
							var Skins []followbot.NamemcSkins
							for i := 2; i < amt+2; i++ {
								request, _ := h2.BuildRequest("https://namemc.com/minecraft-skins"+Path, "GET", "", followbot.CommonHeaders,
									followbot.ReturnCF())

								if _, resp, err := request.Commit(nil); err == nil && resp.StatusCode == 200 {
									resp_body, _ := io.ReadAll(resp.Body)
									Skins = append(Skins, followbot.GetAllSkinsOnPage(string(resp_body))...)
									Path = fmt.Sprintf("?page=%v", i)
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
						Description: "Follow the name you supply from -name",
						Args:        []string{"-name"},
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							Name := StrCmd.String("-name")
							req, _ := http.NewRequest("GET", "https://namemc.com/profile/"+Name, nil)
							req.AddCookie(&http.Cookie{Name: strings.Split(followbot.Cookie, "=")[0], Value: strings.Split(followbot.Cookie, "=")[1]})
							req.AddCookie(&http.Cookie{Name: "session_id", Value: followbot.Session.Key})
							req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
							var redirect string
							Client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { redirect = req.URL.String(); return nil }}
							resp, _ := Client.Do(req)
							if resp.StatusCode == 200 {
								body, _ := io.ReadAll(resp.Body)
								var Target followbot.Target = followbot.Target{
									UUID: strings.Split(strings.Split(string(body), `col-xl" style="font-size: 90%"><samp>`)[1], `</samp></div>`)[0],
									URL:  redirect,
								}
								for _, acc := range followbot.Session.Accounts {
									Target.Send(true)
									if acc.Name != followbot.Session.DisplayName {
										followbot.Swapprofile(acc.URLPath)
									}
								}
							}
						},
					},
					"rl-test": {
						Description: "Scrapes all front page names that are dropping and saves them.",
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

							request, _ := h2.BuildRequest("https://namemc.com/minecraft-names", "GET", "", followbot.CommonHeaders,
								followbot.ReturnCF())

						Exit:
							for {
								var I []Resps
								if _, resp, err := request.Commit(nil); err == nil {
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
											fmt.Print(utils.Logo(fmt.Sprintf("<%v> collected %v name info...\n", time.Now().Format("05.000"), len(Info))))
											request.ChangeURL(fmt.Sprintf("https://namemc.com%v", html.FindOne(doc, "/html/body/main/div/div[5]/nav/ul/li[4]/a").Attr[1].Val))
											if len(I) < 3 {
												break Exit
											}
										}
									} else {
										fmt.Print(utils.Logo("Ratelimited, sleeping for 10 seconds..\n"))
										time.Sleep(10 * time.Second)
									}
								}
							}
							j, _ := json.MarshalIndent(Info, " ", "  ")
							os.Create("namemc.json")
							os.WriteFile("namemc.json", j, 0644)
						},
					},
					"png": {
						Args:        []string{"-name"},
						Description: "Gets the namemc skins of a profile and builds the namemc image from it.",
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}
							var name string = StrCmd.String("-name")

							request, _ := h2.BuildRequest("https://namemc.com/profile/"+name, "GET", "", followbot.CommonHeaders,
								followbot.ReturnCF())
							if _, resp, err := request.Commit(nil); err == nil {
								request.ChangeURL("https://namemc.com" + resp.Header.Get("Location"))
							}

							if _, resp, err := request.Commit(nil); err == nil && resp.StatusCode == 200 {
								if !strings.Contains(resp.Request.URL.String(), "search") {
									if body, err := io.ReadAll(resp.Body); err == nil {
										img := followbot.GenImages(string(body))
										path := "data/skinarts/" + strings.ReplaceAll(uuid.NewString(), "-", "") + "_" + name + ".png"
										if out, err := os.Create(path); err == nil {
											if val, err := png.Decode(bytes.NewBuffer(img)); err == nil {
												if png.Encode(out, val) == nil {
													fmt.Println(utils.Logo(fmt.Sprintf("<%v> Saved %v to %v!", time.Now().Format("05.000"), name, path)))
												}

												// ik its a ugly err chain but what ever.
											} else {
												fmt.Println(utils.Logo(fmt.Sprintf("<%v> [ERROR] %v", time.Now().Format("05.000"), err.Error())))
											}
										} else {
											fmt.Println(utils.Logo(fmt.Sprintf("<%v> [ERROR] %v", time.Now().Format("05.000"), err.Error())))
										}
									} else {
										fmt.Println(utils.Logo(fmt.Sprintf("<%v> [ERROR] %v", time.Now().Format("05.000"), err.Error())))
									}
								} else {
									fmt.Println(utils.Logo(fmt.Sprintf("<%v> [ERROR] %v's Profile Does Not Exist.", time.Now().Format("05.000"), name)))
								}
							} else {
								fmt.Println(utils.Logo(fmt.Sprintf("<%v> [ERROR] %v Unsatisfactory Status Code From namemc.com!", time.Now().Format("05.000"), resp.Status)))
							}
						},
					},
					"3c": {
						Args: []string{"-3n", "-3l"},
						Action: func() {
							isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
							if utils.Con.CF.Tokens == "" {
								return
							}
							values := followbot.ParseAndPerformThreeReqs()
							var Use []followbot.Resps
							switch true {
							case StrCmd.Bool("-3n"):
								for _, username := range values {
									if _, err := strconv.Atoi(username.Name); err == nil {
										Use = append(Use, username)
									}
								}
							case StrCmd.Bool("-3l"):
								for _, username := range values {
									if !isAlpha(username.Name) {
									} else {
										Use = append(Use, username)
									}
								}
							default:
								Use = values
							}
							for _, names := range Use {
								fmt.Println(utils.Logo(fmt.Sprintf(`%v: %v %v | %v`, names.Name, names.Start_Unix, names.End_Unix, names.Index)))
							}
						},
					},
					"mass-follow": {
						Description: "Grab all followers of a namemc profile and send them all follow requests!",
						Args:        []string{"-name", "--single"},
						Action: func() {
							if utils.Con.CF.Tokens == "" {
								return
							}

							type NamemcFollows struct {
								UUID       string
								ProfileURL string
								SkinID     string
								Two        string
								Three      string
							}

							var Data []NamemcFollows

							name := StrCmd.String("-name")
							if name == "" {
								return
							}

							request, _ := h2.BuildRequest("https://namemc.com/profile/"+name, "GET", "", followbot.CommonHeaders,
								followbot.ReturnCF())

							var next string

							if _, resp, err := request.Commit(nil); err == nil {
								request.ChangeURL("https://namemc.com" + resp.Header.Get("Location") + "/followers")
								next = "https://namemc.com" + resp.Header.Get("Location") + "/followers"
							}

							for {
								if _, resp, err := request.Commit(nil); err == nil && resp.StatusCode == 200 {
									body, _ := io.ReadAll(resp.Body)
									if !strings.Contains(string(body), `" rel="next" title="Next Page"`) {
										next = "stop"
									} else {
										d := strings.Split(strings.Split(string(body), `" rel="next" title="Next Page"`)[0], `"`)
										request.ChangeURL(next + "?sort=date:desc&page=" + strings.Split(strings.Split(d[len(d)-1], `?sort=date:desc&amp;page=`)[1], `"`)[0])
									}

									bodys := strings.Split(strings.Split(strings.Split(string(body), "<tbody>")[1], "</tbody>")[0], `<input type="hidden" name="profile" value=`)

									if len(bodys) < 2 {
										break
									}

									for _, data := range bodys {
										if strings.Contains(data, "/profile/") {
											UUID := strings.Split(strings.Split(data, `"`)[1], `">`)[0]
											profileURL := `https://namemc.com/profile/` + strings.Split(strings.Split(data, `<a href="/profile/`)[1], `" `)[0]
											ID := strings.Split(strings.Split(data, `src="https://s.namemc.com/2d/skin/face.png?id=`)[1], `&`)[0]
											twodee := fmt.Sprintf(`https://s.namemc.com/2d/skin/face.png?id=%v&scale=4`, ID)
											threedee := fmt.Sprintf(`https://s.namemc.com/3d/skin/body.png?id=%v&model=classic&width=150&height=200`, ID)
											Data = append(Data, NamemcFollows{
												UUID:       UUID,
												ProfileURL: profileURL,
												SkinID:     ID,
												Two:        twodee,
												Three:      threedee,
											})
										}
									}
									fmt.Println(utils.Logo(fmt.Sprintf(`<!> %v got %v uwus~`, next, len(Data))))
									if next == "stop" {
										break
									}
								}
							}

							fmt.Println()

							Single := StrCmd.Bool("--single")
							var AccName string
							if Single {
								var AccNames []string
								for _, acc := range followbot.Session.Accounts {
									AccNames = append(AccNames, acc.Name)
								}
								fmt.Println(utils.Logo(strings.Join(AccNames, " ")))
								fmt.Print(utils.Logo("What account to single out for follows: "))
								fmt.Scan(&AccName)
								for _, acc := range followbot.Session.Accounts {
									if strings.EqualFold(AccName, acc.Name) {
										followbot.Swapprofile(acc.URLPath)
										break
									}
								}
							}

							for _, names := range Data {
								var Target followbot.Target = followbot.Target{
									UUID: names.UUID,
									URL:  names.ProfileURL,
								}

								if Single {
									Target.Send(false)
									continue
								}

								for _, acc := range followbot.Session.Accounts {
									if Target.Send(true) && acc.Name != followbot.Session.DisplayName {
										followbot.Swapprofile(acc.URLPath)
									} else {
										break
									}
								}
							}
						},
					},
				},
			},
			"key": {
				Description: "Gets your namemc.com key!",
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
			"recover": {
				Description: "Applies recovery emails to invalids.txt accounts.",
				Action: func() {
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
				},
				Args: []string{"--force"},
			},
			"snipe": {
				Description: "Main sniper command, targets only one ign that is passed through with -u",
				Action: func() {
					if len(utils.Bearer.Details) != 0 {

						cl, name, Changed, EmailClaimed := false, StrCmd.String("-u"), false, ""
						var start, end int64 = int64(StrCmd.Int("-start")), int64(StrCmd.Int("-end"))
						if utils.Con.NMC.UseNMC {
							start, end, _, _ = followbot.GetDroptimes(name)
						} else if !utils.Con.NMC.UseNMC || start == 0 || end == 0 {
							if utils.Con.Settings.AskForUnixPrompt {
								fmt.Println(utils.Logo("Timestamp to Unix: [https://www.epochconverter.com/] (make sure to remove the • on the namemc timestamp!)"))
								var collection string
								fmt.Print(utils.Logo("start-end: "))
								fmt.Scan(&collection)
								if data := strings.Split(collection, "-"); len(data) > 1 {
									var start_, end_ = data[0], data[1]
									if val, err := strconv.Atoi(start_); err == nil {
										start = int64(val)
									}
									if val, err := strconv.Atoi(end_); err == nil {
										end = int64(val)
									}
								}
							} else {
								if resp, err := http.Get("https://namemc.info/data/info/" + name); err == nil && resp.StatusCode == 200 {
									var Data utils.NInfo
									json.Unmarshal([]byte(apiGO.ReturnJustString(io.ReadAll(resp.Body))), &Data)
									start = Data.Data.StartDate.Unix()
									end = Data.Data.EndDate.Unix()
								}
							}
						}

						var namemc followbot.NameRequest

						if utils.Con.NMC.UseNMC {
							if start != 0 || end != 0 {
								namemc = followbot.Info(name)
								go func() {
									for !cl || time.Now().Unix() < namemc.Start_Unix {
										time.Sleep(time.Minute)
										namemc = followbot.Info(name)
									}
								}()
							}
						}

						drop := time.Unix(int64(start), 0)

						if time.Now().Before(drop) {
							fmt.Println()
						}

						for time.Now().Before(drop) {
							if utils.Con.NMC.UseNMC {
								if start != 0 || end != 0 {
									fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v | Views - %v | Status - %v                \r", name, time.Until(drop).Round(time.Second), namemc.Searches, namemc.Status))))
								}
							} else {
								fmt.Print(utils.Logo((fmt.Sprintf("[%v] %v                 \r", name, time.Until(drop).Round(time.Second)))))
							}
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
								Spread   time.Duration
							}
							var (
								Payloads []Proxys
							)

							for _, Acc := range utils.Accs["Giftcard"] {
								Payloads = append(Payloads, Proxys{
									Accounts: Acc.Accs,
									Proxy:    Acc.Proxy,
									Spread:   TempCalc(utils.Con.Settings.SleepAmtPerGc, utils.Accamt),
								})
							}

							for _, Acc := range utils.Accs["Microsoft"] {
								Payloads = append(Payloads, Proxys{
									Accounts: Acc.Accs,
									Proxy:    Acc.Proxy,
									Spread:   TempCalc(utils.Con.Settings.SleepAmtPerMfa, utils.Accamt),
								})
							}

							fmt.Println()

							for !cl || !Changed {
								for _, c := range Payloads {
									for _, accs := range c.Accounts {
										go func(Config apiGO.Info, c Proxys) {
											if P, ok, _ := utils.Connect(c.Proxy); ok {
												var wg sync.WaitGroup
												for i := 0; i < map[string]int{"Giftcard": utils.Con.Settings.GC_ReqAmt, "Microsoft": utils.Con.Settings.MFA_ReqAmt}[Config.AccountType]; i++ {
													wg.Add(1)
													go func() {
														if Req := (&apiGO.Details{ResponseDetails: apiGO.SocketSending(P, utils.ReturnPayload(Config.AccountType, Config.Bearer, name)), Bearer: Config.Bearer, Email: Config.Email, Type: Config.AccountType}); Req.ResponseDetails.StatusCode == "200" {
															if utils.Con.SkinChange.Link != "" {
																apiGO.ChangeSkin(apiGO.JsonValue(utils.Con.SkinChange), Config.Bearer)
															}
															NMC := utils.Namemc_key(Config.Bearer)
															if utils.Con.NMC.UseNMC {
																followbot.Claim_NAMEMC(NMC)
																followbot.SendFollow(name)
															}
															EmailClaimed = fmt.Sprint(utils.Success().Apply("✓"), utils.Logo(fmt.Sprintf("%v claimed %v @ %v ↪ ~ %v ~\n", Config.Email, name, time.Now().Format("05.0000"), NMC)))
															cl = true
														} else {
															fmt.Println(utils.Failure().Apply("✗"), utils.Logo(fmt.Sprintf(`<%v> [%v] %v ⑇ %v`, time.Now().Format("05.0000"), Req.ResponseDetails.StatusCode, name, utils.HashEmailClean(Config.Email))))
														}
														wg.Done()
													}()
												}
												wg.Wait()
											}
										}(accs, c)
									}
									time.Sleep(map[bool]time.Duration{true: time.Duration(utils.Con.Settings.Spread) * time.Millisecond, false: c.Spread}[utils.Con.Settings.UseCustomSpread])
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
				Description: "Looks if you have sniped the name you supply via checking your accounts.",
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
			"clear": {
				Description: "Clears the terminal.",
				Action: func() {
					apiGO.Clear()
					fmt.Print(utils.Logo(`/~' _   _ _ |_ | _ 
\_,||_|| | ||_)|(/_

`))
				},
			},
			"yt": {
				Subcommand: map[string]StrCmd.SubCmd{
					"pause": {Action: func() {
						if utils.Con.Settings.Youtube != "" {
							speaker.Lock()
						}
					}},
					"unpause": {Action: func() {
						if utils.Con.Settings.Youtube != "" {
							speaker.Unlock()
						}
					}},
				},
			},
			"spread": {
				Action: func() {
					if utils.Accamt > 0 {
						fmt.Println(utils.Logo(fmt.Sprintf(`<!> GC ↪ %v | MFA ↪ %v`, TempCalc(utils.Con.Settings.SleepAmtPerGc, utils.Accamt), TempCalc(utils.Con.Settings.SleepAmtPerMfa, utils.Accamt))))
					}
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

func TempCalc(interval, accamt int) time.Duration {
	return time.Duration(interval/accamt) * time.Millisecond
}

func Logo() {

	App_.PrintMiddleUncachedToBody(logo)

	time.Sleep(2 * time.Second)

	Center.Clear()

}
