package utils

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

func (s *Config) ToJson() (Data []byte, err error) {
	return json.MarshalIndent(s, "", "  ")
}

func (config *Config) SaveConfig() {
	if Json, err := config.ToJson(); err == nil {
		WriteFile("data/config.json", string(Json))
	}
}

func (s *Config) LoadState() {
	data, err := ReadFile("data/config.json")
	if err != nil {
		s.LoadFromFile()
		s.Settings = AccountSettings{
			GC_ReqAmt:       1,
			MFA_ReqAmt:      1,
			AccountsPerGc:   5,
			AccountsPerMfa:  1,
			SleepAmtPerGc:   15000,
			SleepAmtPerMfa:  10000,
			Spread:          0,
			UseCustomSpread: false,
		}
		s.Bools = Bools{
			FirstUse:           true,
			UseProxyDuringAuth: false,
			DownloadedPW:       false,
			UseWebhook:         false,
		}
		s.SkinChange.Variant = "slim"
		s.SkinChange.Link = "https://textures.minecraft.net/texture/516accb84322ca168a8cd06b4d8cc28e08b31cb0555eee01b64f9175cefe7b75"
		s.Gradient = []Values{{R: "115", G: "52", B: "115"}, {R: "71", G: "33", B: "71"}}
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	jsonFile, err := os.Open("data/config.json")
	if err != nil {
		jsonFile, _ = os.Create("data/config.json")
	}

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func HashEmailClean(email string) string {
	e := strings.Split(email, "@")[0] // stfu
	return e[0:2] + strings.Repeat("⋅", 2) + e[len(e)-5:]
}
