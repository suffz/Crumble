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
		}
		s.SkinChange.Variant = "slim"
		s.SkinChange.Link = "https://textures.minecraft.net/texture/516accb84322ca168a8cd06b4d8cc28e08b31cb0555eee01b64f9175cefe7b75"
		s.Gradient = []Values{{R: "224", G: "217", B: "215"}, {R: "136", G: "147", B: "133"}, {R: "103", G: "126", B: "116"}}
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
	e := strings.Split(email, "@")[0]
	return e[0:len(e)/2] + "____" + "@__" + strings.Split(email, "@")[1][2:]
}
