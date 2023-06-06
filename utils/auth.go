package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"main/apiGO"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/bot/screen"
	"github.com/Tnze/go-mc/bot/world"
	"github.com/Tnze/go-mc/chat"
)

type Proxys_Accs struct {
	Proxy string
	Accs  []apiGO.Info
}

var Accs map[string][]Proxys_Accs = make(map[string][]Proxys_Accs)
var Use_gc, Accamt int
var First_gc bool = true

func AuthAccs() {
	grabDetails()
	if len(Con.Bearers) == 0 {
		fmt.Println(Logo("No Bearers have been found, please check your details."))
		return
	} else {
		checkifValid()
		for _, Accs := range Con.Bearers {
			if Accs.NameChange {
				Bearer.Details = append(Bearer.Details, apiGO.Info{
					Bearer:      Accs.Bearer,
					AccountType: Accs.Type,
					Email:       Accs.Email,
					Password:    Accs.Password,
					Info:        apiGO.UserINFO(Accs.Info),
				})
			}
		}
	}
}

func grabDetails() {

	var AccountsVer []string
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		fmt.Println(Logo("Unable to continue, you have no Accounts added."))
		return
	}
	P := Payload(AccountsVer)
	CheckDupes(AccountsVer)
	if Con.Bearers == nil {
		P_Auth(P, false)
	} else if len(Con.Bearers) < len(AccountsVer) {
		var auth []string
		check := make(map[string]bool)
		for _, Acc := range Con.Bearers {
			check[Acc.Email+":"+Acc.Password] = true
		}
		for _, Accs := range AccountsVer {
			if !check[Accs] {
				auth = append(auth, Accs)
			}
		}
		P_Auth(Payload(auth), false)
	} else if len(AccountsVer) < len(Con.Bearers) {
		var New []Bearers
		for _, Accs := range AccountsVer {
			for _, num := range Con.Bearers {
				if Accs == num.Email+":"+num.Password {
					New = append(New, num)
					break
				}
			}
		}
		Con.Bearers = New
	}

	Con.SaveConfig()
	Con.LoadState()
}

func checkifValid() {
	var reAuth []string
	var wgs sync.WaitGroup
	for _, Accs := range Con.Bearers {
		if time.Now().Unix() > Accs.AuthedAt+Accs.AuthInterval {
			reAuth = append(reAuth, Accs.Email+":"+Accs.Password)
		} else {
			if Accs.NameChange {
				wgs.Add(1)
				go func(Accs Bearers) {
					f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
					f.Header.Set("Authorization", "Bearer "+Accs.Bearer)
					if j, err := http.DefaultClient.Do(f); err == nil {
						if j.StatusCode == 401 {
							reAuth = append(reAuth, Accs.Email+":"+Accs.Password)
						}
					}
					wgs.Done()
				}(Accs)
			}
		}
	}
	wgs.Wait()
	if len(reAuth) != 0 {
		P_Auth(Payload(reAuth), true)
	}
	Con.SaveConfig()
	Con.LoadState()
}

// _diamondburned_#4507 thanks to them for the epic example below.

func CheckDupes(strs []string) []string {
	dedup := strs[:0] // re-use the backing array
	track := make(map[string]bool, len(strs))

	for _, str := range strs {
		if track[str] {
			continue
		}
		dedup = append(dedup, str)
		track[str] = true
	}

	return dedup
}

func CheckAccs() {
	for {
		time.Sleep(10 * time.Second)
		var reauth []string
		for _, acc := range Con.Bearers {
			if time.Now().Unix() > acc.AuthedAt+acc.AuthInterval && acc.NameChange {
				reauth = append(reauth, acc.Email+":"+acc.Password)
			}
		}
		if len(reauth) > 0 {
			P_Auth(Payload(reauth), true)
		}
		Con.SaveConfig()
		Con.LoadState()
	}
}

type Payload_auth struct {
	Proxy    string
	Accounts []string
}

func Payload(accounts []string) (Data []Payload_auth) {
	var use_proxy, ug int
	var f bool = true
	for _, bearer := range accounts {
		if use_proxy >= len(Proxy.Proxys) && len(Proxy.Proxys) < len(Bearer.Details) {
			break
		}

		if f {
			Data = append(Data, Payload_auth{Proxy: Proxy.Proxys[use_proxy]})
			f = false
			use_proxy++
		}
		if len(Data[ug].Accounts) != 3 {
			Data[ug].Accounts = append(Data[ug].Accounts, bearer)
		} else {
			ug++
			Data = append(Data, Payload_auth{Proxy: Proxy.Proxys[use_proxy], Accounts: []string{bearer}})
			use_proxy++
		}
	}
	return
}

var (
	client        *bot.Client
	player        *basic.Player
	chatHandler   *msg.Manager
	worldManager  *world.World
	screenManager *screen.Manager
)

func Namemc_key(bearer string) string {

	var Info_Acc apiGO.UserINFO

	if len(Proxy.Proxys) > 0 {
		p := Proxy.CompRand()
		info := strings.Split(p, ":")
		var ip, port, user, pass string
		switch len(info) {
		case 2:
			ip = info[0]
			port = info[1]
		case 4:
			ip = info[0]
			port = info[1]
			user = info[2]
			pass = info[3]
		}
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

func NotWhiteSpace(str string) bool {
	for _, c := range str {
		if !unicode.IsSpace(c) {
			return true
		}
	}
	return false
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

func ReturnJustString(data []byte, err error) string {
	return string(data)
}

func P_Auth(P []Payload_auth, reauth bool) {
	var wg sync.WaitGroup
	var Invalids []string
	var invalidproxys []string
	for _, acc_1 := range P {
		for _, p := range acc_1.Accounts {
			if acc := strings.Split(p, ":"); len(acc) > 1 {
				if len(Proxy.Proxys) > 0 && Con.Bools.UseProxyDuringAuth {
					wg.Add(1)
					go func(proxy string, acc []string) {
						ip, port, user, pass := "", "", "", ""
						switch data := strings.Split(proxy, ":"); len(data) {
						case 2:
							ip = data[0]
							port = data[1]
						case 4:
							ip = data[0]
							port = data[1]
							user = data[2]
							pass = data[3]
						}
						var Authed bool
						go func() {
							for !Authed {
								time.Sleep(80 * time.Second)
								if !Authed {
									New := Proxy.CompRand()
									fmt.Println(Logo(fmt.Sprintf("<%v> %v under the ip %v has timed out, suspected dead proxy. reauthing under %v now..", time.Now().Format("05.000"), HashEmailClean(acc[0]), ip, strings.Split(New, ":")[0])))

									invalidproxys = append(invalidproxys, fmt.Sprintf("%v:%v:%v:%v", ip, port, user, pass))

									switch data := strings.Split(New, ":"); len(data) {
									case 2:
										ip = data[0]
										port = data[1]
									case 4:
										ip = data[0]
										port = data[1]
										user = data[2]
										pass = data[3]
									}
									info := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass})
									if info.Error != "" {
										Authed = true
										fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", HashEmailClean(info.Email), info.Error)))
										Invalids = append(Invalids, acc[0]+":"+acc[1])
									} else if info.Bearer != "" {
										if IsChangeable(proxy, info.Bearer) {
											Authed = true
											fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashEmailClean(info.Email))))
											if reauth {
												for point, bf := range Con.Bearers {
													if strings.EqualFold(bf.Email, info.Email) {
														Con.Bearers[point] = Bearers{
															Bearer:       info.Bearer,
															NameChange:   true,
															Type:         info.AccountType,
															Password:     info.Password,
															Email:        info.Email,
															AuthedAt:     time.Now().Unix(),
															AuthInterval: 54000,
															Info: UserINFO{
																ID:   info.Info.ID,
																Name: info.Info.Name,
															},
														}
														break
													}
												}
												for i, Bearers := range Bearer.Details {
													if strings.EqualFold(Bearers.Email, info.Email) {
														Bearer.Details[i] = info
														break
													}
												}
												var Found bool
											E1:
												for i, accs := range Accs["Giftcard"] {
													for e, b := range accs.Accs {
														if strings.EqualFold(b.Email, info.Email) {
															Accs["Giftcard"][i].Accs[e] = info
															Found = true
															break E1
														}
													}
												}
												if !Found {
												E2:
													for i, accs := range Accs["Microsoft"] {
														for e, b := range accs.Accs {
															if strings.EqualFold(b.Email, info.Email) {
																Accs["Microsoft"][i].Accs[e] = info
																Found = true
																break E2
															}
														}
													}
												}
											} else {
												Con.Bearers = append(Con.Bearers, Bearers{
													Bearer:       info.Bearer,
													AuthInterval: 54000,
													AuthedAt:     time.Now().Unix(),
													Type:         info.AccountType,
													Email:        info.Email,
													Password:     info.Password,
													NameChange:   true,
													Info: UserINFO{
														ID:   info.Info.ID,
														Name: info.Info.Name,
													},
												})
											}
										} else {
											fmt.Println(Logo(fmt.Sprintf("Account %v cannot name change.. %v", acc[0], info.Info.Name)))
											for point, bf := range Con.Bearers {
												if strings.EqualFold(bf.Email, info.Email) {
													Con.Bearers[point] = Bearers{
														Type:         info.AccountType,
														Bearer:       info.Bearer,
														NameChange:   false,
														Password:     info.Password,
														Email:        info.Email,
														AuthedAt:     time.Now().Unix(),
														AuthInterval: 54000,
														Info:         UserINFO(info.Info),
													}
													break
												}
											}
											for i, Bearers := range Bearer.Details {
												if strings.EqualFold(Bearers.Email, info.Email) {
													Bearer.Details[i] = info
													break
												}
											}
											var Found bool
										E13:
											for i, accs := range Accs["Giftcard"] {
												for e, b := range accs.Accs {
													if strings.EqualFold(b.Email, info.Email) {
														Accs["Giftcard"][i].Accs[e] = info
														Found = true
														break E13
													}
												}
											}
											if !Found {
											E23:
												for i, accs := range Accs["Microsoft"] {
													for e, b := range accs.Accs {
														if strings.EqualFold(b.Email, info.Email) {
															Accs["Microsoft"][i].Accs[e] = info
															Found = true
															break E23
														}
													}
												}
											}
											Invalids = append(Invalids, acc[0]+":"+acc[1])
										}
									}
									wg.Done()
								}
							}
						}()
						info := apiGO.MS_authentication(acc[0], acc[1], &apiGO.ProxyMS{IP: ip, Port: port, User: user, Password: pass})
						if info.Error != "" {
							Authed = true
							fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", HashEmailClean(info.Email), info.Error)))
							Invalids = append(Invalids, acc[0]+":"+acc[1])
						} else if info.Bearer != "" {
							if IsChangeable(proxy, info.Bearer) {
								Authed = true
								fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashEmailClean(info.Email))))
								if reauth {
									for point, bf := range Con.Bearers {
										if strings.EqualFold(bf.Email, info.Email) {
											Con.Bearers[point] = Bearers{
												Bearer:       info.Bearer,
												NameChange:   true,
												Type:         info.AccountType,
												Password:     info.Password,
												Email:        info.Email,
												AuthedAt:     time.Now().Unix(),
												AuthInterval: 54000,
												Info: UserINFO{
													ID:   info.Info.ID,
													Name: info.Info.Name,
												},
											}
											break
										}
									}
									for i, Bearers := range Bearer.Details {
										if strings.EqualFold(Bearers.Email, info.Email) {
											Bearer.Details[i] = info
											break
										}
									}
									var Found bool
								E1:
									for i, accs := range Accs["Giftcard"] {
										for e, b := range accs.Accs {
											if strings.EqualFold(b.Email, info.Email) {
												Accs["Giftcard"][i].Accs[e] = info
												Found = true
												break E1
											}
										}
									}
									if !Found {
									E2:
										for i, accs := range Accs["Microsoft"] {
											for e, b := range accs.Accs {
												if strings.EqualFold(b.Email, info.Email) {
													Accs["Microsoft"][i].Accs[e] = info
													Found = true
													break E2
												}
											}
										}
									}
								} else {
									Con.Bearers = append(Con.Bearers, Bearers{
										Bearer:       info.Bearer,
										AuthInterval: 54000,
										AuthedAt:     time.Now().Unix(),
										Type:         info.AccountType,
										Email:        info.Email,
										Password:     info.Password,
										NameChange:   true,
										Info: UserINFO{
											ID:   info.Info.ID,
											Name: info.Info.Name,
										},
									})
								}
							} else {
								fmt.Println(Logo(fmt.Sprintf("Account %v cannot name change.. %v", acc[0], info.Info.Name)))
								for point, bf := range Con.Bearers {
									if strings.EqualFold(bf.Email, info.Email) {
										Con.Bearers[point] = Bearers{
											Type:         info.AccountType,
											Bearer:       info.Bearer,
											NameChange:   false,
											Password:     info.Password,
											Email:        info.Email,
											AuthedAt:     time.Now().Unix(),
											AuthInterval: 54000,
											Info:         UserINFO(info.Info),
										}
										break
									}
								}
								for i, Bearers := range Bearer.Details {
									if strings.EqualFold(Bearers.Email, info.Email) {
										Bearer.Details[i] = info
										break
									}
								}
								var Found bool
							E13:
								for i, accs := range Accs["Giftcard"] {
									for e, b := range accs.Accs {
										if strings.EqualFold(b.Email, info.Email) {
											Accs["Giftcard"][i].Accs[e] = info
											Found = true
											break E13
										}
									}
								}
								if !Found {
								E23:
									for i, accs := range Accs["Microsoft"] {
										for e, b := range accs.Accs {
											if strings.EqualFold(b.Email, info.Email) {
												Accs["Microsoft"][i].Accs[e] = info
												Found = true
												break E23
											}
										}
									}
								}
								Invalids = append(Invalids, acc[0]+":"+acc[1])
							}
						}
						wg.Done()
					}(acc_1.Proxy, acc)
				} else {
					switch info := apiGO.MS_authentication(acc[0], acc[1], nil); true {
					case info.Error != "":
						fmt.Println(Logo(fmt.Sprintf("Account %v came up Invalid: %v", HashEmailClean(info.Email), info.Error)))
						Invalids = append(Invalids, acc[0]+":"+acc[1])
					case info.Bearer != "" && IsChangeable("", info.Bearer):
						fmt.Println(Logo(fmt.Sprintf("[%v] Succesfully authed %v", time.Now().Format("15:04:05.0000"), HashEmailClean(info.Email))))
						if reauth {
							for point, bf := range Con.Bearers {
								if strings.EqualFold(bf.Email, info.Email) {
									Con.Bearers[point] = Bearers{
										Bearer:       info.Bearer,
										NameChange:   true,
										Type:         info.AccountType,
										Password:     info.Password,
										Email:        info.Email,
										AuthedAt:     time.Now().Unix(),
										AuthInterval: 54000,
										Info: UserINFO{
											ID:   info.Info.ID,
											Name: info.Info.Name,
										},
									}
									break
								}
							}
							for i, Bearers := range Bearer.Details {
								if strings.EqualFold(Bearers.Email, info.Email) {
									Bearer.Details[i] = info
									break
								}
							}
							var Found bool
						E1:
							for i, accs := range Accs["Giftcard"] {
								for e, b := range accs.Accs {
									if strings.EqualFold(b.Email, info.Email) {
										Accs["Giftcard"][i].Accs[e] = info
										Found = true
										break E1
									}
								}
							}
							if !Found {
							E2:
								for i, accs := range Accs["Microsoft"] {
									for e, b := range accs.Accs {
										if strings.EqualFold(b.Email, info.Email) {
											Accs["Microsoft"][i].Accs[e] = info
											Found = true
											break E2
										}
									}
								}
							}
						} else {
							Con.Bearers = append(Con.Bearers, Bearers{
								Bearer:       info.Bearer,
								AuthInterval: 54000,
								AuthedAt:     time.Now().Unix(),
								Type:         info.AccountType,
								Email:        info.Email,
								Password:     info.Password,
								NameChange:   true,
							})
						}
					default:
						fmt.Println(Logo(fmt.Sprintf("Account %v Bearer is nil or it cannot name change.. [%v]", HashEmailClean(acc[0]), acc[1])))
						Invalids = append(Invalids, acc[0]+":"+acc[1])
					}
				}
			}
		}
	}
	wg.Wait()
	if len(Invalids) != 0 {
		if _, err := os.Stat("invalids.txt"); os.IsNotExist(err) {
			os.Create("invalids.txt")
		}

		os.WriteFile("invalids.txt", []byte(strings.Join(Invalids, "\n")), 0644)

		scanner := bufio.NewScanner(strings.NewReader(strings.Join(Invalids, "\n")))
		for scanner.Scan() {
			for i, acc := range Con.Bearers {
				if strings.EqualFold(acc.Email, strings.Split(scanner.Text(), ":")[0]) {
					Con.Bearers[i].NameChange = false
					Con.SaveConfig()
					Con.LoadState()
					break
				}
			}
		}
	}
	if len(invalidproxys) != 0 {
		if body, err := os.ReadFile("proxys.txt"); err == nil {
			strings.ReplaceAll(string(body), strings.Join(invalidproxys, "\n"), "")
			fmt.Println(strings.ReplaceAll(string(body), strings.Join(invalidproxys, "\n"), ""))
		}
	}
}

func IsChangeable(proxy, bearer string) bool {
	data := strings.Split(proxy, ":")
	var Client http.Client
	switch len(data) {
	case 4:
		Client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", User: url.UserPassword(data[2], data[3]), Host: data[0] + ":" + data[1]})}}
	case 2:
		Client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: data[0] + ":" + data[1]})}}
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
