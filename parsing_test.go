package torrent

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCreateAddr(t *testing.T) {
	addr := createAddr("http", "duckduckgo.com", "html", "Алексей Пехов Страж")
	answer := "http://duckduckgo.com/html?ia=web&num=50&q=%D0%90%D0%BB%D0%B5%D0%BA%D1%81%D0%B5%D0%B9+%D0%9F%D0%B5%D1%85%D0%BE%D0%B2+%D0%A1%D1%82%D1%80%D0%B0%D0%B6+torrent"
	if addr != answer {
		t.Errorf("результат функции %s,\nправильный результат: %s", addr, answer)
	}
}

func TestSearchResult(t *testing.T) {
	addr := createAddr("http", "duckduckgo.com", "html", "Алексей Пехов Страж")
	responseBody, err := getSiteBody(addr)
	if err != nil {
		t.Error(err.Error())
	}
	ioutil.WriteFile("file.html", responseBody, 0666)
	results := getSearchResult(responseBody)
	if len(results) == 0 {
		t.Error("нет результатов поиска")
	}
	fmt.Println(strings.Join(results, "\n"))
}

func TestParsingSites(t *testing.T) {
	addr := createAddr("http", "duckduckgo.com", "html", "Капитал Маркс аудиокнига")
	responseBody, err := getSiteBody(addr)
	if err != nil {
		t.Error(err.Error())
	}
	//ioutil.WriteFile("file.html", responseBody, 0666)
	results := getSearchResult(responseBody)
	if len(results) == 0 {
		t.Error("нет результатов поиска")
	}
	torrents := parsingSites(results)
	if len(torrents) == 0 {
		t.Error("нет найденных торрент-файлов")
	}
	fmt.Println(strings.Join(torrents, "\n"))
}
