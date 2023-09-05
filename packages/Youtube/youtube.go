package Youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/packages/utils"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Video(U Youtube, Video bool, Quality string) YTRequest {
	ID := U.ID
	Body := fmt.Sprintf(`
	{
	  "context": {
		"client": {
		  "clientName": "WEB",
		  "clientVersion": "2.20230615.02.01"
		}
	  },
	  "videoId": "%v"
	}
	`, ID)
	req, _ := http.NewRequest("POST", "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false", bytes.NewReader([]byte(Body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(Body)))
	resp, e := http.DefaultClient.Do(req)
	if e == nil {
		defer resp.Body.Close()
		var getOpts YT
		get_id, _ := io.ReadAll(resp.Body)
		json.Unmarshal(get_id, &getOpts)

		for _, T := range getOpts.StreamingData.AdaptiveFormats {
			if Video {
				if T.QualityLabel == Quality {

					DL := Youtube{
						ID:      ID,
						FullURL: "https://www.youtube.com/watch?v=" + ID,
						Title:   getOpts.VideoDetails.Title,
						Length:  T.ContentLength,
						Index:   getOpts.VideoDetails.ShortDescription,
					}

					return DL.getReq(T.URL, T.Sig, T.ContentLength)
				}
			} else {
				if T.AudioQuality == Quality {

					DL := Youtube{
						ID:      ID,
						FullURL: "https://www.youtube.com/watch?v=" + ID,
						Title:   getOpts.VideoDetails.Title,
						Length:  T.ContentLength,
						Index:   getOpts.VideoDetails.ShortDescription,
					}

					return DL.getReq(T.URL, T.Sig, T.ContentLength)
				}
			}
		}
		for _, T := range getOpts.StreamingData.Formats {
			if Video {
				if T.QualityLabel == Quality {

					DL := Youtube{
						ID:      ID,
						FullURL: "https://www.youtube.com/watch?v=" + ID,
						Title:   getOpts.VideoDetails.Title,
						Length:  resp.Header.Get("Content-Length"),
						Index:   getOpts.VideoDetails.ShortDescription,
					}

					return DL.getReq(T.URL, T.Sig, resp.Header.Get("Content-Length"))
				}
			} else {
				if T.AudioQuality == Quality {

					DL := Youtube{
						ID:      ID,
						FullURL: "https://www.youtube.com/watch?v=" + ID,
						Title:   getOpts.VideoDetails.Title,
						Length:  resp.Header.Get("Content-Length"),
						Index:   getOpts.VideoDetails.ShortDescription,
					}

					return DL.getReq(T.URL, T.Sig, resp.Header.Get("Content-Length"))
				}
			}
		}
	}
	return YTRequest{}
}

func Playlist(url string) (IDs []Youtube) {
	req, _ := http.NewRequest("GET", url, nil)
	rr, _ := http.DefaultClient.Do(req)
	aa, _ := io.ReadAll(rr.Body)
	var DD utils.YTPageConfig
	json.Unmarshal([]byte(strings.Split(strings.Split(string(aa), `var ytInitialData =`)[1], `;</script>`)[0]), &DD)
	for _, data := range DD.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		for _, yt := range data.TabRenderer.Content.SectionListRenderer.Contents {
			for _, pagedata := range yt.ItemSectionRenderer.Contents {
				for _, data := range pagedata.PlaylistVideoListRenderer.Contents {
					if data.PlaylistVideoRenderer.IsPlayable && data.PlaylistVideoRenderer.VideoID != "" {
						IDs = append(IDs, Youtube{
							ID:           data.PlaylistVideoRenderer.VideoID,
							Title:        getTitle(data.PlaylistVideoRenderer.Title.Runs),
							FullURL:      "https://www.youtube.com/watch?v=" + data.PlaylistVideoRenderer.VideoID,
							Continuation: data.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token,
							Length:       data.PlaylistVideoRenderer.LengthText.SimpleText,
							Index:        data.PlaylistVideoRenderer.Index.SimpleText,
							Info:         getTitle(data.PlaylistVideoRenderer.VideoInfo.Runs),
						})
					} else {
						if data.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token != "" {
							IDs[len(IDs)-1].Continuation = data.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
						}
					}
				}
			}
		}
	}
	return
}

func (DL *YTRequest) Download() ([]byte, time.Duration, error) {
	st := time.Now()
	ct, err := strconv.Atoi(DL.ContentLength)
	if err != nil {
		return []byte{}, 0, err
	}
	var Range string
	if DL.Sig {
		Range = fmt.Sprintf(`&range=0-%v`, ct)
	}
	req, _ := http.NewRequest("GET", DL.URL+Range, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Range", fmt.Sprintf(fmt.Sprintf("bytes=0-%v", ct)))
	req.Header.Add("Referer", DL.URL)
	r, _ := http.DefaultClient.Do(req)
	Vid, err := io.ReadAll(r.Body)
	return Vid, time.Since(st), err
}

func ReturnJustString(data []byte, err error) string {
	return string(data)
}

func (DL *Youtube) getReq(URL, SIG, CT string) YTRequest {
	if URL != "" {
		if int_value, err := strconv.Atoi(CT); err == nil {
			return YTRequest{
				VideoID:       DL.ID,
				ContentLength: strconv.Itoa(int_value),
				URL:           URL,
				Config:        *DL,
			}
		}
	}
	if SIG != "" {
		if CL, Url, err := getSigUrlAndToken(SIG, DL.ID); err == nil {
			return YTRequest{
				VideoID:       DL.ID,
				ContentLength: strconv.Itoa(CL),
				URL:           Url,
				Config:        *DL,
				Sig:           true,
			}
		}
	}
	return YTRequest{}
}

func getSigUrlAndToken(SIG, VideoID string) (int, string, error) {
	pars, err := url.ParseQuery(SIG)
	if err != nil {
		return 0, "", err
	}

	u, err := url.Parse(pars.Get("url"))
	if err != nil {
		return 0, "", err
	}

	S, err := url.PathUnescape(pars.Get("s"))
	if err != nil {
		return 0, "", err
	}

	a, err := decrypt([]byte(S), VideoID)
	if err != nil {
		return 0, "", err
	}

	S = string(a)

	// decode S and get the token.
	q := u.Query()

	config, _ := getPlayerConfig(VideoID)
	q.Add(pars.Get("sp"), S)

	vals, err := decryptNParam(config, q)
	if err != nil {
		return 0, "", err
	}

	u.RawQuery = vals.Encode()

	URL := u.String()

	// perform the request to get the content length, as sometimes with the sig value requests it doesnt give one in the json itself.
	resp, err := http.Get(URL)
	if err != nil {
		return 0, "", err
	}

	ContentL, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return 0, "", err
	}
	return ContentL, URL, err
}

func Write(Vid []byte, filename string) *os.File {
	files, _ := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0644)
	io.Copy(files, bytes.NewBuffer(Vid))
	return files
}

func getTitle(Runs []struct {
	Text string "json:\"text\""
}) (fill string) {
	for _, t := range Runs {
		fill += t.Text
	}
	return
}

func (YT YTRequest) Play(inp, out string, body []byte) {

	if ffmpeg, ok := checkFFM(); ok {
		file := Write(body, inp)
		file.Close()

		cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", inp, "-vn", out)
		cmd.Run()
		os.Remove(inp)

		file, err := os.Open(out)
		if err != nil {
			return
		}

		streamer, format, err := mp3.Decode(file)
		if err != nil {
			panic(err)
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))
		<-done
		speaker.Clear()
		speaker.Close()
	}
}

func checkFFM() (string, bool) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		if runtime.GOOS == "windows" {
			if err := exec.Command("winget", "install", "ffmpeg").Run(); err != nil {
				return ffmpeg, false
			} else {
				ffmpeg, _ = exec.LookPath("ffmpeg")
			}
		}
		return ffmpeg, true
	}
	return ffmpeg, true
}
