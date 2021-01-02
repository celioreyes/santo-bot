package howlongtobeat

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	hltbURL           = "https://howlongtobeat.com"
	hltbSearchResults = "/search_results"
)

type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HowLongToBeat -- wrapper around HowLongToBeat.com for getting game information
type Service struct {
	baseURL string
	client  httpclient
}

// GameInfo contains a game title and the various times to beat it
type GameInfo struct {
	Title string
	Times []GameTime
}

// GameTime contains the type of time and the value for a game
type GameTime struct {
	Type  string
	Value string
}

func New(client httpclient) *Service {
	return &Service{
		baseURL: hltbURL,
		client:  client,
	}
}

func (s *Service) createGameSearchRequest(game string) (req *http.Request, err error) {

	// TODO: make this cleaner and extensible
	// default all values except game
	form := url.Values{}
	form.Add("queryString", game)
	form.Add("t", "games")
	form.Add("sorthead", "popular")
	form.Add("sortd", "Normal Order")
	form.Add("plat", "")
	form.Add("length_type", "main")
	form.Add("length_min", "")
	form.Add("length_max", "")
	form.Add("detail", "")
	form.Add("randomize", "")

	path := fmt.Sprintf("%s%s?page=1", s.baseURL, hltbSearchResults)
	req, err = http.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return
}

func (s *Service) SearchGame(game string) (games []GameInfo, err error) {

	req, err := s.createGameSearchRequest(game)
	if err != nil {
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// read body
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	doc.Find("ul li").Each(func(i int, item *goquery.Selection) {
		// both time and type of time are under the same selector
		// query for them and match every pair together
		var game GameInfo

		game.Title = item.Find(".search_list_details h3 a").Text()

		timeInfo := item.Find(".search_list_tidbit")
		for i := 0; i < len(timeInfo.Nodes); i = i + 2 {
			timeType := timeInfo.Eq(i).Text()
			timeValue := timeInfo.Eq(i + 1).Text()

			gameTime := GameTime{
				Type:  timeType,
				Value: timeValue,
			}

			game.Times = append(game.Times, gameTime)
		}

		games = append(games, game)
	})

	return
}
