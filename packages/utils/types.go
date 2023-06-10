package utils

import (
	"crypto/x509"
	"main/packages/webhook"
	"time"

	"main/packages/apiGO"
)

var (
	Roots  *x509.CertPool = x509.NewCertPool()
	Con    Config
	Proxy  apiGO.Proxys
	Bearer apiGO.MCbearers
	RGB    []string
)

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
	WebhookURL string          `json:"webhook_url"`
	Webhook    webhook.Web     `json:"webhook_json"`
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
	UseMethod                        bool `json:"use_method_rlbypass"`
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
