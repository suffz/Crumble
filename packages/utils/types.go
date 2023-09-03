package utils

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"main/packages/apiGO"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/screen"
	"github.com/Tnze/go-mc/bot/world"
)

var (
	Roots                   *x509.CertPool = x509.NewCertPool()
	Con                     Config
	Proxy                   apiGO.Proxys
	Bearer                  apiGO.MCbearers
	RGB                     []string
	First_mfa               bool = true
	First_gc                bool = true
	Use_gc, Use_mfa, Accamt int
	Accs                    map[string][]Proxys_Accs = make(map[string][]Proxys_Accs)
	client                  *bot.Client
	player                  *basic.Player
	chatHandler             *msg.Manager
	worldManager            *world.World
	screenManager           *screen.Manager
	letterRunes             = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	ProxyByte               = []byte(`
	-- GlobalSign Root R2, valid until Dec 15, 2021
	-----BEGIN CERTIFICATE-----
	MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
	A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
	Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
	MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
	A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
	hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
	v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
	eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
	tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
	C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
	zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
	mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
	V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
	bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
	3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
	J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
	291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
	ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
	AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
	TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
	-----END CERTIFICATE-----`)
)

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

type Proxys_Accs struct {
	Proxy string
	Accs  []apiGO.Info
}

type NameMCInfo struct {
	Action string     `json:"action"`
	Desc   string     `json:"desc"`
	Code   string     `json:"code"`
	Data   NameMCData `json:"data"`
}
type NameMCData struct {
	Status    string    `json:"status"`
	Searches  string    `json:"searches"`
	Begin     int       `json:"begin"`
	End       int       `json:"end"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Headurl   string    `json:"headurl"`
}

type NameMCHead struct {
	Bodyurl string `json:"bodyurl"`
	Headurl string `json:"headurl"`
	ID      string `json:"id"`
}

type Names struct {
	Name  string
	Taken bool
}

type Proxies struct {
	IP, Port, User, Password string
}

type Status struct {
	Data struct {
		Status string `json:"status"`
	} `json:"details"`
}

type CF struct {
	Tokens   string `json:"tokens"`
	GennedAT int64  `json:"unix_of_creation"`
}

type Config struct {
	Gradient   []Values        `json:"gradient"`
	NMC        Namemc_Data     `json:"namemc_settings"`
	Settings   AccountSettings `json:"settings"`
	Bools      Bools           `json:"sniper_config"`
	SkinChange Skin            `json:"skin_config"`
	CF         CF              `json:"cf_tokens"`
	Bearers    []Bearers       `json:"Bearers"`
	Recovery   []Succesful     `json:"recovery"`
}

type Namemc_Data struct {
	UseNMC          bool       `json:"usenamemc_fordroptime_andautofollow"`
	Display         string     `json:"name_to_use_for_follows"`
	Key             string     `json:"namemc_email:pass"`
	NamemcLoginData NMC        `json:"namemc_login_data"`
	P               []Profiles `json:"genned_profiles"`
}

type Profiles struct {
	Session_ID string `json:"session_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type NMC struct {
	Token      string `json:"token"`
	LastAuthed int64  `json:"last_unix_auth_timestamp"`
}

type Bearers struct {
	Bearer               string   `json:"Bearer"`
	Email                string   `json:"Email"`
	Password             string   `json:"Password"`
	AuthInterval         int64    `json:"AuthInterval"`
	AuthedAt             int64    `json:"AuthedAt"`
	Type                 string   `json:"Type"`
	NameChange           bool     `json:"NameChange"`
	Info                 UserINFO `json:"Info"`
	NOT_ENTITLED_CHECKED bool     `json:"checked_entitled"`
}

type Succesful struct {
	Email     string
	Recovery  string
	Code_Used string
}
type Data struct {
	Info []Succesful
}

type Refresh struct {
	Time_since_last_gen int64 `json:"last_entitled_prevention"`
}

type AccountSettings struct {
	AskForUnixPrompt bool  `json:"ask_for_unix_prompt"`
	AccountsPerGc    int   `json:"accounts_per_gc_proxy"`
	AccountsPerMfa   int   `json:"accounts_per_mfa_proxy"`
	GC_ReqAmt        int   `json:"amt_reqs_per_gc_acc"`
	MFA_ReqAmt       int   `json:"amt_reqs_per_mfa_acc"`
	SleepAmtPerGc    int   `json:"sleep_for_gc"`
	SleepAmtPerMfa   int   `json:"sleep_for_mfa"`
	UseCustomSpread  bool  `json:"use_own_spread_value"`
	Spread           int64 `json:"spread_ms"`
}

type Bools struct {
	UseCF                            bool `json:"use_cf_token"`
	UseProxyDuringAuth               bool `json:"useproxysduringauth"`
	UseWebhook                       bool `json:"sendpersonalwhonsnipe"`
	FirstUse                         bool `json:"firstuse_IGNORETHIS"`
	DownloadedPW                     bool `json:"pwinstalled_IGNORETHIS"`
	ApplyNewRecoveryToExistingEmails bool `json:"applynewemailstoexistingrecoveryemails"`
}

type Values struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

type Info struct {
	Bearer       string
	RefreshToken string
	AccessToken  string
	Expires      int
	AccountType  string
	Email        string
	Password     string
	Requests     int
	Info         UserINFO `json:"Info"`
	Error        string
}

type UserINFO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Skin struct {
	Link    string `json:"url"`
	Variant string `json:"variant"`
}

type Payload_auth struct {
	Proxy    string
	Accounts []string
}
