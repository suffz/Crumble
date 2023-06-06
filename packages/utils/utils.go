package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/iskaa02/qalam/gradient"
)

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
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

func GetHeadUrl(name string) (string, string) {
	if resp, err := http.Get("https://namemc.info/data/namemc/head/" + name); err == nil {
		res, _ := io.ReadAll(resp.Body)
		var Data NameMCHead
		json.Unmarshal(res, &Data)
		return Data.Headurl, Data.Bodyurl
	}
	return "", ""
}

func GetDroptimes(name string) (int64, int64, string, string) {
	if resp, err := http.Get("https://namemc.info/data/info/" + name); err == nil {
		res, _ := io.ReadAll(resp.Body)
		var Data NameMCInfo
		json.Unmarshal(res, &Data)
		return int64(Data.Data.Begin), int64(Data.Data.End), Data.Data.Status, Data.Data.Searches
	}
	return 0, 0, "", ""
}
