package Youtube

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/dop251/goja"
)

const (
	Video2160p  = "2160p"
	Video1440p  = "1440p"
	Video1080p  = "1080p"
	Video720p   = "720p"
	Video480p   = "480p"
	Video360p   = "360p"
	Video240p   = "240p"
	Video144p   = "144p"
	AudioMedium = "AUDIO_QUALITY_MEDIUM"
	AudioLow    = "AUDIO_QUALITY_LOW"
)

type YTVids struct {
	RBody io.ReadCloser
	Body  []byte
	Index int
}

type YTRequest struct {
	ContentLength string
	URL           string
	VideoID       string
	Sig           bool
	Config        Youtube
}

type Youtube struct {
	ID           string
	FullURL      string
	Title        string
	Length       string
	Index        string
	Info         string
	Continuation string
}

type YT struct {
	StreamingData StreamingData `json:"streamingData"`
	VideoDetails  Details       `json:"videoDetails"`
}
type Details struct {
	VideoID          string   `json:"videoId"`
	Title            string   `json:"title"`
	LengthSeconds    string   `json:"lengthSeconds"`
	Keywords         []string `json:"keywords"`
	ChannelID        string   `json:"channelId"`
	IsOwnerViewing   bool     `json:"isOwnerViewing"`
	ShortDescription string   `json:"shortDescription"`
	IsCrawlable      bool     `json:"isCrawlable"`
	Thumbnail        struct {
		Thumbnails []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	AllowRatings      bool   `json:"allowRatings"`
	ViewCount         string `json:"viewCount"`
	Author            string `json:"author"`
	IsPrivate         bool   `json:"isPrivate"`
	IsUnpluggedCorpus bool   `json:"isUnpluggedCorpus"`
	IsLiveContent     bool   `json:"isLiveContent"`
}
type Formats struct {
	Itag             int    `json:"itag"`
	URL              string `json:"url"`
	MimeType         string `json:"mimeType"`
	Bitrate          int    `json:"bitrate"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	LastModified     string `json:"lastModified"`
	Quality          string `json:"quality"`
	Fps              int    `json:"fps"`
	QualityLabel     string `json:"qualityLabel"`
	ProjectionType   string `json:"projectionType"`
	AudioQuality     string `json:"audioQuality"`
	ApproxDurationMs string `json:"approxDurationMs"`
	AudioSampleRate  string `json:"audioSampleRate"`
	AudioChannels    int    `json:"audioChannels"`
	Sig              string `json:"signatureCipher"`
}
type InitRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type IndexRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type ColorInfo struct {
	Primaries               string `json:"primaries"`
	TransferCharacteristics string `json:"transferCharacteristics"`
	MatrixCoefficients      string `json:"matrixCoefficients"`
}
type AdaptiveFormats struct {
	Itag             int        `json:"itag"`
	URL              string     `json:"url"`
	MimeType         string     `json:"mimeType"`
	Bitrate          int        `json:"bitrate"`
	Width            int        `json:"width,omitempty"`
	Height           int        `json:"height,omitempty"`
	InitRange        InitRange  `json:"initRange"`
	IndexRange       IndexRange `json:"indexRange"`
	LastModified     string     `json:"lastModified"`
	ContentLength    string     `json:"contentLength"`
	Quality          string     `json:"quality"`
	Fps              int        `json:"fps,omitempty"`
	QualityLabel     string     `json:"qualityLabel,omitempty"`
	ProjectionType   string     `json:"projectionType"`
	AverageBitrate   int        `json:"averageBitrate"`
	ApproxDurationMs string     `json:"approxDurationMs"`
	ColorInfo        ColorInfo  `json:"colorInfo,omitempty"`
	HighReplication  bool       `json:"highReplication,omitempty"`
	AudioQuality     string     `json:"audioQuality,omitempty"`
	AudioSampleRate  string     `json:"audioSampleRate,omitempty"`
	AudioChannels    int        `json:"audioChannels,omitempty"`
	LoudnessDb       float64    `json:"loudnessDb,omitempty"`
	Sig              string     `json:"signatureCipher"`
}
type StreamingData struct {
	ExpiresInSeconds string            `json:"expiresInSeconds"`
	Formats          []Formats         `json:"formats"`
	AdaptiveFormats  []AdaptiveFormats `json:"adaptiveFormats"`
}

type DecipherOperation func([]byte) []byte

var Y = regexp.MustCompile(`^.*(youtu.be\/|v\/|e\/|u\/\w+\/|embed\/|v=)([^#\&\?]*).*`)

func YoutubeURL(URL string) string {
	YT_ := Y.FindAllStringSubmatch(URL, -1)
	if len(YT_) > 0 {
		if len(YT_[0]) > 2 {
			return YT_[0][2]
		}
	}
	return "Unknown"
}

func decrypt(cyphertext []byte, id string) ([]byte, error) {
	operations, err := parseDecipherOps(cyphertext, id)
	if err != nil {
		return nil, err
	}

	// apply operations
	bs := []byte(cyphertext)
	for _, op := range operations {
		bs = op(bs)
	}

	return bs, nil
}

const (
	jsvarStr   = "[a-zA-Z_\\$][a-zA-Z_0-9]*"
	reverseStr = ":function\\(a\\)\\{" +
		"(?:return )?a\\.reverse\\(\\)" +
		"\\}"
	spliceStr = ":function\\(a,b\\)\\{" +
		"a\\.splice\\(0,b\\)" +
		"\\}"
	swapStr = ":function\\(a,b\\)\\{" +
		"var c=a\\[0\\];a\\[0\\]=a\\[b(?:%a\\.length)?\\];a\\[b(?:%a\\.length)?\\]=c(?:;return a)?" +
		"\\}"
)

var (
	nFunctionNameRegexp = regexp.MustCompile("\\.get\\(\"n\"\\)\\)&&\\(b=([a-zA-Z0-9$]{0,3})\\[(\\d+)\\](.+)\\|\\|([a-zA-Z0-9]{0,3})")
	actionsObjRegexp    = regexp.MustCompile(fmt.Sprintf(
		"var (%s)=\\{((?:(?:%s%s|%s%s|%s%s),?\\n?)+)\\};", jsvarStr, jsvarStr, swapStr, jsvarStr, spliceStr, jsvarStr, reverseStr))
	actionsFuncRegexp = regexp.MustCompile(fmt.Sprintf(
		"function(?: %s)?\\(a\\)\\{"+
			"a=a\\.split\\(\"\"\\);\\s*"+
			"((?:(?:a=)?%s\\.%s\\(a,\\d+\\);)+)"+
			"return a\\.join\\(\"\"\\)"+
			"\\}", jsvarStr, jsvarStr, jsvarStr))
	reverseRegexp = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, reverseStr))
	spliceRegexp  = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, spliceStr))
	swapRegexp    = regexp.MustCompile(fmt.Sprintf("(?m)(?:^|,)(%s)%s", jsvarStr, swapStr))
)

func decryptNParam(config []byte, query url.Values) (url.Values, error) {
	// decrypt n-parameter
	nSig := query.Get("v")
	if nSig != "" {
		nDecoded, err := decodeNsig(config, nSig)
		if err != nil {
			return nil, fmt.Errorf("unable to decode nSig: %w", err)
		}
		query.Set("v", nDecoded)
	}

	return query, nil
}

func parseDecipherOps(config []byte, id string) (operations []DecipherOperation, err error) {
	config, _ = getPlayerConfig(id)
	objResult := actionsObjRegexp.FindSubmatch(config)
	funcResult := actionsFuncRegexp.FindSubmatch(config)
	if len(objResult) < 3 || len(funcResult) < 2 {
		return nil, fmt.Errorf("error parsing signature tokens (#obj=%d, #func=%d)", len(objResult), len(funcResult))
	}
	obj := objResult[1]
	objBody := objResult[2]
	funcBody := funcResult[1]
	var reverseKey, spliceKey, swapKey string
	if result := reverseRegexp.FindSubmatch(objBody); len(result) > 1 {
		reverseKey = string(result[1])
	}
	if result := spliceRegexp.FindSubmatch(objBody); len(result) > 1 {
		spliceKey = string(result[1])
	}
	if result := swapRegexp.FindSubmatch(objBody); len(result) > 1 {
		swapKey = string(result[1])
	}
	regex, err := regexp.Compile(fmt.Sprintf("(?:a=)?%s\\.(%s|%s|%s)\\(a,(\\d+)\\)", regexp.QuoteMeta(string(obj)), regexp.QuoteMeta(reverseKey), regexp.QuoteMeta(spliceKey), regexp.QuoteMeta(swapKey)))
	if err != nil {
		return nil, err
	}
	var ops []DecipherOperation
	for _, s := range regex.FindAllSubmatch(funcBody, -1) {
		switch string(s[1]) {
		case reverseKey:
			ops = append(ops, reverseFunc)
		case swapKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, newSwapFunc(arg))
		case spliceKey:
			arg, _ := strconv.Atoi(string(s[2]))
			ops = append(ops, newSpliceFunc(arg))
		}
	}
	return ops, nil
}

var basejsPattern = regexp.MustCompile(`(/s/player/\w+/player_ias.vflset/\w+/base.js)`)

func getPlayerConfig(videoID string) ([]byte, error) {
	embedURL := fmt.Sprintf("https://youtube.com/embed/%s?hl=en", videoID)
	req, _ := http.NewRequest("GET", embedURL, nil)
	req.Header.Set("Origin", "https://youtube.com")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	resp, _ := http.DefaultClient.Do(req)
	embedBody, _ := io.ReadAll(resp.Body)
	playerPath := string(basejsPattern.Find(embedBody))
	if playerPath == "" {
		return nil, errors.New("unable to find basejs URL in playerConfig")
	}
	reqa, _ := http.NewRequest("GET", "https://youtube.com"+playerPath, nil)
	reqa.Header.Set("Origin", "https://youtube.com")
	reqa.Header.Set("Sec-Fetch-Mode", "navigate")
	respa, _ := http.DefaultClient.Do(reqa)
	re, _ := io.ReadAll(respa.Body)
	return re, nil
}

func reverseFunc(bs []byte) []byte {
	l, r := 0, len(bs)-1
	for l < r {
		bs[l], bs[r] = bs[r], bs[l]
		l++
		r--
	}
	return bs
}

func newSwapFunc(arg int) DecipherOperation {
	return func(bs []byte) []byte {
		pos := arg % len(bs)
		bs[0], bs[pos] = bs[pos], bs[0]
		return bs
	}
}

func newSpliceFunc(pos int) DecipherOperation {
	return func(bs []byte) []byte {
		return bs[pos:]
	}
}

func getNFunction(config []byte) (string, error) {
	nameResult := nFunctionNameRegexp.FindSubmatch(config)
	if len(nameResult) == 0 {
		return "", errors.New("unable to extract n-function name")
	}

	var name string
	if idx, _ := strconv.Atoi(string(nameResult[2])); idx == 0 {
		name = string(nameResult[4])
	} else {
		name = string(nameResult[1])
	}

	return extraFunction(config, name)

}

func decodeNsig(config []byte, encoded string) (string, error) {
	fBody, err := getNFunction(config)
	if err != nil {
		return "", err
	}

	return evalJavascript(fBody, encoded)
}

func evalJavascript(jsFunction, arg string) (string, error) {
	const myName = "myFunction"

	vm := goja.New()
	_, err := vm.RunString(myName + "=" + jsFunction)
	if err != nil {
		return "", err
	}

	var output func(string) string
	err = vm.ExportTo(vm.Get(myName), &output)
	if err != nil {
		return "", err
	}

	return output(arg), nil
}

func extraFunction(config []byte, name string) (string, error) {
	// find the beginning of the function
	def := []byte(name + "=function(")
	start := bytes.Index(config, def)
	if start < 1 {
		return "", fmt.Errorf("unable to extract n-function body: looking for '%s'", def)
	}

	// start after the first curly bracket
	pos := start + bytes.IndexByte(config[start:], '{') + 1

	var strChar byte

	// find the bracket closing the function
	for brackets := 1; brackets > 0; pos++ {
		b := config[pos]
		switch b {
		case '{':
			if strChar == 0 {
				brackets++
			}
		case '}':
			if strChar == 0 {
				brackets--
			}
		case '`', '"', '\'':
			if config[pos-1] == '\\' && config[pos-2] != '\\' {
				continue
			}
			if strChar == 0 {
				strChar = b
			} else if strChar == b {
				strChar = 0
			}
		}
	}

	return string(config[start:pos]), nil
}
