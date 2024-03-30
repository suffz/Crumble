package followbot

import (
	"image"
	"time"

	"github.com/faiface/beep"
	"github.com/suffz/Youtube"
)

var (
	Session Namemc
	Cookie  string
)

type YTVIDEO struct {
	Vid       Youtube.Youtube
	PlayedAlr bool
	S         beep.StreamSeekCloser
	F         beep.Format
	Streamer  Youtube.YTRequest
}

type Images struct {
	Image image.Image
	Url   string
	Row   int
}

type Skins struct {
	DownloadURL string `json:"download,omitempty"`
	URL         string `json:"url,omitempty"`
	ID          string `json:"id,omitempty"`
	ChangedAt   string `json:"changedat,omitempty"`
	Head        string `json:"headurl,omitempty"`
	Body        string `json:"bodyurl,omitempty"`
}

type Namemc struct {
	Key         string
	DisplayName string
	Accounts    []UUIDS
}

type UUIDS struct {
	Name    string
	URLPath string
	IconPNG string
}

type NameRequest struct {
	Status     string     `json:"status,omitempty"`
	Searches   string     `json:"searches,omitempty"`
	Start_Unix int64      `json:"begin,omitempty"`
	End_Unix   int64      `json:"end,omitempty"`
	Start      *time.Time `json:"start_date,omitempty"`
	End        *time.Time `json:"end_date,omitempty"`
	HeadURL    string     `json:"headurl,omitempty"`
	BodyUrl    string     `json:"bodyurl,omitempty"`
	Error      string     `json:"error,omitempty"`
}

type NamemcSkins struct {
	Emoji          string `json:"emoji"`
	NamemcUsername string `json:"owner"`
	Number         string `json:"number"`
	Stars          string `json:"stars"`
	Time           string `json:"time"`
	BodyURL        string `json:"bodyurl"`
	HeadURL        string `json:"headurl"`
	SkinDownload   string `json:"skindownload"`
}

type Target struct {
	UUID string
	URL  string
	Hits int
}

type Resps struct {
	Name       string     `json:"name,omitempty" bson:"name"`
	Start      *time.Time `json:"start_date,omitempty"`
	End        *time.Time `json:"end_date,omitempty"`
	Start_Unix int64      `json:"start_unix,omitempty"`
	End_Unix   int64      `json:"end_unix,omitempty"`
	CachedAt   string     `json:"cachedat,omitempty" bson:"cachedat"`
	Searches   string     `json:"searches,omitempty" bson:"searches"`
	Index      int
}
