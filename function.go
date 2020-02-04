// Package p contains an HTTP Cloud Function.
package p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getBooksLink(source []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	links := []string{}
	doc.Find(".bookitem").Each(func(i int, s *goquery.Selection) {
		links = append(links, s.Find("a").AttrOr("href", ""))
	})
	return links, nil
}

func getBook(url string, ch chan string) {
	source, err := download(url)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(source))
	if err != nil {
		return
	}
	title := doc.Find(".bwname .mar0pad0").Text()

	date := doc.Find("[itemprop='publish_date']").Parent().Text()
	date = strings.NewReplacer(" ", "", "\n", "", "發售日：", "").Replace(date)

	ch <- fmt.Sprintf("%s: %s", title, date)
}

func getVol(title string) int {
	re := regexp.MustCompile(`\d+`)
	m := re.FindAll([]byte(title), -1)
	v, _ := strconv.Atoi(string(m[0]))
	return v
}

type Input struct {
	Sn int `json:"sn"`
}

func (data *Input) GetUrl() string {
	sn := data.Sn
	if sn == 0 {
		sn = 927
	}
	return fmt.Sprintf("https://www.bookwalker.com.tw/search?series=%d&order=sell", sn)
}

func Fetch(w http.ResponseWriter, r *http.Request) {
	var in Input
	s, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(s, &in); err != nil {
	}

	source, err := download(in.GetUrl())
	if err != nil {
		fmt.Println(err)
		return
	}

	links, err := getBooksLink(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan string)
	for _, v := range links {
		go getBook(v, ch)
	}

	var title, date string
	var vol int
	for i := 0; i < len(links); i++ {
		select {
		case msg := <-ch:
			data := strings.Split(msg, ": ")
			v := getVol(data[0])
			if title == "" || data[1] > date || (data[1] == date && v > vol) {
				title = data[0]
				date = data[1]
				vol = v
			}
		}
	}
	fmt.Fprintf(w, "%s %s", title, date)
}
