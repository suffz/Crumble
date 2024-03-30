package apiGO

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"main/packages/h2"

	tls2 "github.com/bogdanfinn/utls"
)

func (accountBearer MCbearers) CreatePayloads(name string) (Data Payload) {
	for _, bearer := range accountBearer.Details {
		if bearer.AccountType == "Giftcard" {
			Data.Payload = append(Data.Payload, fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%s\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %s\r\n\r\n"+string([]byte(`{"profileName":"`+name+`"}`))+"\r\n", strconv.Itoa(len(string([]byte(`{"profileName":"`+name+`"}`)))), bearer.Bearer))
		} else {
			Data.Payload = append(Data.Payload, "PUT /minecraft/profile/name/"+name+" HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: MCSN/1.0\r\nContent-Length:0\r\nAuthorization: bearer "+bearer.Bearer+"\r\n\r\n")
		}
	}

	return
}

func Sleep(dropTime int64, delay float64) {
	time.Sleep(time.Until(time.Unix(dropTime, 0).Add(time.Millisecond * time.Duration(0-delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))
}

func GetConfig(owo []byte) (Config map[string]interface{}) {
	json.Unmarshal(owo, &Config)
	return
}

func Sum(array []float64) (sum float64) {
	for _, ammount := range array {
		sum = sum + ammount
	}

	return
}

func CheckChange(bearer string, PS *h2.ProxyAuth) bool {
	var Conn h2.Conn
	var err error
	if PS != nil {
		Conn, err = (&h2.Client{Config: h2.GetDefaultConfig()}).Connect("https://api.minecraftservices.com/minecraft/profile/namechange", h2.ReqConfig{ID: 1, BuildID: tls2.HelloChrome_120_PQ, DataBodyMaxLength: 1609382, Proxy: PS})
	} else {
		Conn, err = (&h2.Client{Config: h2.GetDefaultConfig()}).Connect("https://api.minecraftservices.com/minecraft/profile/namechange", h2.ReqConfig{ID: 1, BuildID: tls2.HelloChrome_120_PQ, DataBodyMaxLength: 1609382})
	}
	if err != nil {
		return false
	} else {
		Conn.AddHeader("Authorization", "Bearer "+bearer)
		if resp, err := Conn.Do("GET", "", "", nil); err == nil {
			switch resp.Status {
			case "200":
				var Data struct {
					NC bool `json:"nameChangeAllowed"`
				}
				json.Unmarshal(resp.Data, &Data)
				return Data.NC
			case "404":
				return true
			default:
				return false
			}
		} else {
			return false
		}
	}
}

func SocketSending(conn net.Conn, payload string) Resp {
	fmt.Fprintln(conn, payload)
	sendTime := time.Now()
	recvd := make([]byte, 4069)
	conn.Read(recvd)
	return Resp{
		RecvAt:     time.Now(),
		StatusCode: string(recvd[9:12]),
		SentAt:     sendTime,
		Body:       string(recvd),
	}
}

func ChangeSkin(body []byte, bearer string) (Req *http.Response, err error) {
	resp, err := http.NewRequest("POST", "https://api.minecraftservices.com/minecraft/profile/skins", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp.Header.Set("Authorization", "bearer "+bearer)

	Req, err = http.DefaultClient.Do(resp)
	if err != nil {
		return nil, err
	}

	return
}

var (
	UrlPost           = regexp.MustCompile(`urlPost:'(.+?)'`)
	Value             = regexp.MustCompile(`value="(.*?)"`)
	User_Authenticate = `{"Properties": {"AuthMethod": "RPS", "SiteName": "user.auth.xboxlive.com", "RpsTicket": "%v"}, "RelyingParty": "http://auth.xboxlive.com", "TokenType": "JWT"}`
	Xsts_Authenticate = `{"Properties": {"SandboxId": "RETAIL", "UserTokens": ["%v"]}, "RelyingParty": "rp://api.minecraftservices.com/", "TokenType": "JWT"}`
	Login_With_Xbox   = `{"identityToken" : "XBL3.0 x=%v;%v", "ensureLegacyEnabled" : true}`
)

type ProxyMS struct {
	IP       string
	Port     string
	User     string
	Password string
}

func MS_authentication(e, p string, PS *ProxyMS) (returnDetails Info) {
	redirect, bearerMS := "", mojangData{}
	if jar, err := cookiejar.New(nil); err == nil {
		Client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { redirect = req.URL.String(); return nil }, Transport: &http.Transport{TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateFreelyAsClient}}, Jar: jar}
		if resp, err := Client.Get("https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en"); err == nil {
			jar.Cookies(resp.Request.URL)
			body := ReturnJustString(io.ReadAll(resp.Body))
			var v, u string
			if Data := Value.FindAllStringSubmatch(string(body), -1); len(Data) > 0 && len(Data[0]) > 0 {
				v = Data[0][1]
			} else {
				fmt.Println("entered reauth for "+e, "Error occured From VALUE find")
				returnDetails = MS_authentication(e, p, PS)
			}
			if Data := UrlPost.FindAllStringSubmatch(string(body), -1); len(Data) > 0 && len(Data[0]) > 0 {
				u = Data[0][1]
			} else {
				fmt.Println("entered reauth for "+e, "Error occured from getting UrlPost")
				returnDetails = MS_authentication(e, p, PS)
			}
			Client.Post(u, "application/x-www-form-urlencoded", bytes.NewReader([]byte(fmt.Sprintf("login=%v&loginfmt=%v&passwd=%v&PPFT=%v", url.QueryEscape(e), url.QueryEscape(e), url.QueryEscape(p), v))))
			if strings.Contains(redirect, "access_token") {
				if d, err := Client.Post("https://user.auth.xboxlive.com/user/authenticate", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(User_Authenticate, strings.Split(strings.Split(strings.Split(redirect, "#")[1], "&")[0], "=")[1])))); err == nil {
					var TOKEN string
					Body := ReturnJustString(io.ReadAll(d.Body))
					if strings.Contains(Body, `"Token":"`) {
						TOKEN = strings.Split(strings.Split(Body, `"Token":"`)[1], `"`)[0]
						if x, err := Client.Post("https://xsts.auth.xboxlive.com/xsts/authorize", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(Xsts_Authenticate, TOKEN)))); err == nil {
							x_data := ReturnJustString(io.ReadAll(x.Body))
							if x.StatusCode == 401 && strings.Contains(string(x_data), "XErr") {
								switch true {
								case strings.Contains(string(x_data), "2148916238"):
									returnDetails = Info{Email: e, Password: p, Error: "Account belongs to someone under 18 and needs to be added to a family"}
									return
								case strings.Contains(string(x_data), "2148916233"):
									returnDetails = Info{Email: e, Password: p, Error: "Account has no Xbox account, you must sign up for one first"}
									return
								}
							} else {
								if PS != nil {
									var UHS, TOKEN string
									if strings.Contains(string(x_data), `"uhs":"`) {
										UHS = strings.Split(strings.Split(string(x_data), `"uhs":"`)[1], `"`)[0]
									} else {
										returnDetails = MS_authentication(e, p, PS)
										return
									}
									if strings.Contains(string(x_data), `"Token":"`) {
										TOKEN = strings.Split(strings.Split(string(x_data), `"Token":"`)[1], `"`)[0]
									} else {
										returnDetails = MS_authentication(e, p, PS)
										return
									}
									req, _ := http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBufferString(fmt.Sprintf(Login_With_Xbox, UHS, TOKEN)))
									if resp, err := (&http.Client{Timeout: time.Second * 60, Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: PS.IP + ":" + PS.Port, User: url.UserPassword(PS.User, PS.Password)})}}).Do(req); err == nil {
										switch resp.StatusCode {
										case 200:
											body, _ := io.ReadAll(resp.Body)
											json.Unmarshal([]byte(ReturnJustString(io.ReadAll(bytes.NewBuffer(body)))), &bearerMS)
											Info_MCINFO, AccT := ReturnAll(bearerMS.Bearer_MS, PS)
											returnDetails = Info{Info: Info_MCINFO, Bearer: bearerMS.Bearer_MS, AccessToken: strings.Split(strings.Split(redirect, "access_token=")[1], "&")[0], RefreshToken: strings.Split(strings.Split(redirect, "refresh_token=")[1], "&")[0], Expires: ReturnJustInt(strconv.Atoi(strings.Split(strings.Split(redirect, "expires_in=")[1], "&")[0])), Email: e, Password: p, AccountType: AccT}
										case 429:
											fmt.Printf("[%v] %v Has been ratelimited, sleeping for 1 minute..\n", resp.Status, e)
											time.Sleep(60 * time.Second)
											returnDetails = MS_authentication(e, p, PS)
										default:
											body, _ := io.ReadAll(resp.Body)
											returnDetails = Info{Email: e, Password: p, Error: fmt.Sprintf("[%v] Unknown status code while authenticating.\n%v", resp.Status, string(body))}
										}
									} else {
										returnDetails = Info{Email: e, Password: p, Error: err.Error()}
									}
								} else {
									var UHS, TOKEN string
									if strings.Contains(string(x_data), `"uhs":"`) {
										UHS = strings.Split(strings.Split(string(x_data), `"uhs":"`)[1], `"`)[0]
									} else {
										returnDetails = MS_authentication(e, p, PS)
										return
									}
									if strings.Contains(string(x_data), `"Token":"`) {
										TOKEN = strings.Split(strings.Split(string(x_data), `"Token":"`)[1], `"`)[0]
									} else {
										returnDetails = MS_authentication(e, p, PS)
										return
									}
									if resp, err := http.Post("https://api.minecraftservices.com/authentication/login_with_xbox", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(Login_With_Xbox, UHS, TOKEN)))); err == nil {
										switch resp.StatusCode {
										case 200:
											body, _ := io.ReadAll(resp.Body)
											json.Unmarshal([]byte(ReturnJustString(io.ReadAll(bytes.NewBuffer(body)))), &bearerMS)
											Info_MCINFO, AccT := ReturnAll(bearerMS.Bearer_MS, nil)
											returnDetails = Info{Info: Info_MCINFO, Bearer: bearerMS.Bearer_MS, AccessToken: strings.Split(strings.Split(redirect, "access_token=")[1], "&")[0], RefreshToken: strings.Split(strings.Split(redirect, "refresh_token=")[1], "&")[0], Expires: ReturnJustInt(strconv.Atoi(strings.Split(strings.Split(redirect, "expires_in=")[1], "&")[0])), Email: e, Password: p, AccountType: AccT}
										case 429:
											fmt.Printf("[%v] %v Has been ratelimited, sleeping for 1 minute..\n", resp.Status, e)
											time.Sleep(60 * time.Second)
											returnDetails = MS_authentication(e, p, PS)
										default:
											body, _ := io.ReadAll(resp.Body)
											returnDetails = Info{Email: e, Password: p, Error: fmt.Sprintf("[%v] Unknown status code while authenticating.\n%v", resp.Status, string(body))}
										}
									} else {
										returnDetails = Info{Email: e, Password: p, Error: err.Error()}
									}
								}
							}
						} else {
							returnDetails = Info{Email: e, Password: p, Error: err.Error()}
						}
					} else {
						returnDetails = MS_authentication(e, p, PS)
					}
				} else {
					returnDetails = Info{Email: e, Password: p, Error: err.Error()}
				}
			} else {
				returnDetails = Info{Email: e, Password: p, Error: "Unable to authorize, access_token missing from request."}
			}
		}
	} else {
		returnDetails = MS_authentication(e, p, PS)
	}
	return
}

func ReturnJustInt(i int, e error) int {
	return i
}

func ReturnJustString(data []byte, err error) string {
	return string(data)
}

type Headers struct {
	Name  string
	Value string
}

func GetReqStartedAndBuildHeaders(method, url string, body io.Reader, headers ...Headers) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	for _, name := range headers {
		req.Header.Add(name.Name, name.Value)
	}
	return req
}

type UserINFO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ReturnAll(bearer string, PS *ProxyMS) (Data UserINFO, Accounttype string) {
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

func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
