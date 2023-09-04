package apiGO

import (
	"crypto/tls"
	"time"
)

var (
	ProxyByte = []byte(`
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

type Name struct {
	Names string  `json:"name"`
	Drop  float64 `json:"droptime"`
}

type Proxys struct {
	Proxys   []string
	Used     map[string]bool
	Accounts []Info
	Conn     *tls.Conn
}

type Resp struct {
	SentAt     time.Time
	RecvAt     time.Time
	StatusCode string
	Body       string
}

type Payload struct {
	Payload []string
	Conns   []*tls.Conn

	Start int64 `json:"unix"`
	End   int64 `json:"unix_end"`
}

type ServerInfo struct {
	Webhook string
	SkinUrl string
}

type mojangData struct {
	Bearer_MS string `json:"access_token"`
	Expires   int    `json:"expires_in"`
}

type MCbearers struct {
	Details []Info
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
	Error        string
	Info         UserINFO
}

type Config struct {
	ChangeSkinLink    string `json:"ChangeSkinLink"`
	ChangeskinOnSnipe bool   `json:"ChangeskinOnSnipe"`
	GcReq             int    `json:"GcReq"`
	MFAReq            int    `json:"MFAReq"`
	ManualBearer      bool   `json:"ManualBearer"`
	SpreadPerReq      int    `json:"SpreadPerReq"`
	SendMCSNAd        bool   `json:"SendMCSN_MsgUponBMJoin"`
	AlwaysUseUnix     bool   `json:"forceuseunix"`

	Bearers []Bearers `json:"Bearers"`
	Logs    []Logs    `json:"logs"`
}

type Logs struct {
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Send     time.Time `json:"send"`
	Recv     time.Time `json:"recv"`
	Success  bool      `json:"success"`
	Name     string    `json:"name"`
}

type Bearers struct {
	Bearer       string `json:"Bearer"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	AuthInterval int64  `json:"AuthInterval"`
	AuthedAt     int64  `json:"AuthedAt"`
	Type         string `json:"Type"`
	NameChange   bool   `json:"NameChange"`
}

type Bux2 struct {
	Action string        `json:"action"`
	Desc   string        `json:"desc"`
	Code   string        `json:"code"`
	ID     string        `json:"id"`
	Error  string        `json:"error"`
	Data   []NameRequest `json:"data"`
}

type Bux struct {
	Action string      `json:"action"`
	Desc   string      `json:"desc"`
	Code   string      `json:"code"`
	ID     string      `json:"id"`
	Error  string      `json:"error"`
	Data   NameRequest `json:"data"`
}

type NameRequest struct {
	Status   string `json:"status,omitempty"`
	Searches string `json:"searches,omitempty"`
	Start    int64  `json:"begin,omitempty"`
	End      int64  `json:"end,omitempty"`
	HeadURL  string `json:"headurl,omitempty"`
	Error    string `json:"error,omitempty"`
}

type ReqConfig struct {
	Name     string
	Delay    float64
	Droptime int64
	Proxys   Proxys
	Bearers  MCbearers
	UseUnix  bool

	Proxy bool
}

type SentRequests struct {
	Requests []Details
}

type Details struct {
	ResponseDetails Resp
	Bearer          string
	Email           string
	Password        string
	Type            string
}
