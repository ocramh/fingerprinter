package lookup

// import (
// 	"encoding/json"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strconv"

// 	"github.com/go-fingerprint/fingerprint"
// 	"github.com/go-fingerprint/gochroma"
// 	"github.com/xlab/api"
// )

// const (
// 	testApiKey = "S5pMOkfMeW"
// )

// var (
// 	duration int
// )

// func init() {
// 	flag.IntVar(&duration, "d", 180, `duration of input audio stream
// 		in seconds`)
// }

// func main() {
// 	flag.Parse()

// 	if flag.NArg() < 1 {
// 		println("Usage: lookup [-d=duration] <file>")
// 		os.Exit(0)
// 	}

// 	// Create new fingerprint calculator
// 	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
// 	defer fpcalc.Close()

// 	f, err := os.Open(flag.Arg(0))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Get fingerprint as base64-encoded string
// 	fprint, err := fpcalc.Fingerprint(
// 		fingerprint.RawInfo{
// 			Src:        f,
// 			Channels:   2,
// 			Rate:       44100,
// 			MaxSeconds: 180,
// 		})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("file fingerprint is: ", fprint)

// 	// Determine if our fingerprint corresponds to any song in AcoustId
// 	// database
// 	i, err := lookupSong(fprint, testApiKey, duration)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(i)
// }

// func lookupSong(fprint, apikey string, duration int) (i *songInfo, err error) {
// 	svc, err := api.New("http://api.acoustid.org/v2")

// 	if err != nil {
// 		return
// 	}

// 	v := url.Values{}

// 	v.Set("client", apikey)
// 	v.Set("fingerprint", "d54bc360-e1b9-4e55-a3e8-b887658a1c4a")
// 	v.Set("duration", strconv.Itoa(duration))
// 	v.Set("meta", "recordings releasegroups compress")

// 	req, _ := svc.Request(api.GET, "/lookup", v)

// 	var cli http.Client

// 	resp, err := cli.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, _ := ioutil.ReadAll(resp.Body)

// 	var m struct {
// 		Results []struct {
// 			Recordings []struct {
// 				Duration      int
// 				Releasegroups []struct {
// 					Title string
// 				}
// 				Title   string
// 				Artists []struct {
// 					Name string
// 				}
// 			}
// 		}
// 	}

// 	if err := json.Unmarshal(body, &m); err != nil {
// 		return nil, err
// 	}

// 	if len(m.Results) == 0 {
// 		return nil, errors.New("No results found")
// 	}

// 	return &songInfo{
// 		m.Results[0].Recordings[0].Title,
// 		m.Results[0].Recordings[0].Artists[0].Name,
// 		m.Results[0].Recordings[0].Releasegroups[0].Title,
// 	}, nil
// }

// type songInfo struct {
// 	Name   string
// 	Artist string
// 	Album  string
// }

// func (s *songInfo) String() string {
// 	return fmt.Sprintf("Track:\t%v\nArtist:\t%v\nAlbum:\t%v", s.Name, s.Artist,
// 		s.Album)
// }
