package torrent

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

//createAddr - формирует адрес с параметрами get-запроса
func createAddr(scheme, host, path, search string) string {
	query := url.Values{}
	if search != "" {
		query.Set("q", search+" torrent")
		query.Add("ia", "web")
		query.Add("num", "50")

	}
	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	return u.String()
}

func getSiteBody(addr string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36 OPR/70.0.3728.59 (Edition beta)")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	return body, nil
}

//"#r1-0 > div > h2 > a.result__a"
func getSearchResult(body []byte) []string {
	reader := bytes.NewReader(body)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}
	result := make([]string, 0)
	document.Find("a.result__snippet").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if len(href) > 0 {
			result = append(result, href)
		}
	})
	return result
}

func parsingSites(sites []string) []string {
	var wg sync.WaitGroup
	result := make([]string, 0)
	wg.Add(len(sites))
	for _, site := range sites {
		go func(addr string) {
			defer wg.Done()
			body, err := getSiteBody(addr)
			if err != nil {
				return
			}
			torrent, err := searchTorrentFile(body)
			if err != nil {
				return
			}
			result = append(result, torrent)
		}(site)
	}
	wg.Wait()
	return result
}
