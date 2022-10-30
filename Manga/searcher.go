package Manga

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type searchPage struct {
	Name string
	Url  string
}

type Manga struct {
	Title       string
	Description string
	Chapter     []Chapter
}

type Chapter struct {
	Vol string
	Url string
}

const (
	search_url = "https://manganato.com/search/story/"
)

func (s Manga) GetMangaInfo() {
	fmt.Println("\n" + s.Title + "\n" + s.Description + "\n")
	for _, i := range s.Chapter {
		fmt.Println(i.Vol + " : " + i.Url)
	}
}

func replace_space(str string) string {
	result := strings.ReplaceAll(str, " ", "-")
	return result
}

func MangaSearch(name string) ([]searchPage, error) {
	c := colly.NewCollector()

	finSearch := []searchPage{}

	c.OnHTML("div.search-story-item", func(e *colly.HTMLElement) {
		s := searchPage{
			Name: e.ChildAttr("a", "title"),
			Url:  e.ChildAttr("a", "href"),
		}
		finSearch = append(finSearch, s)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(search_url + replace_space(name))
	c.Wait()

	if len(finSearch) == 0 {
		return finSearch, errors.New("NO RESULT")
	}

	return finSearch, nil
}

func MangaPage(url string) Manga {
	c := colly.NewCollector()
	m := Manga{}

	c.OnHTML("h1", func(h *colly.HTMLElement) {
		m.Title = h.Text
	})

	c.OnHTML("div.panel-story-info-description", func(h *colly.HTMLElement) {
		m.Description = h.Text
	})

	c.OnHTML("li.a-h", func(e *colly.HTMLElement) {
		cha := Chapter{}

		cha.Vol = e.ChildText("a")
		cha.Url = e.ChildAttr("a", "href")
		m.Chapter = append(m.Chapter, cha)
	})

	c.Visit(url)
	c.Wait()
	return m
}

func GetMangaPicLink(url string) []string {
	c := colly.NewCollector()
	var urlList []string
	c.OnHTML("div.container-chapter-reader", func(h *colly.HTMLElement) {
		urlList = h.ChildAttrs("img", "src")
	})

	c.Visit(url)
	c.Wait()
	return urlList
}

func DownloadFromPage(ur []string, chap int, name string) error {
	client := &http.Client{}
	dirName, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	os.Mkdir(dirName+"/Documents/"+name, os.ModePerm)
	os.Mkdir(dirName+"/Documents/"+name+"/"+strconv.Itoa(chap), os.ModePerm)
	for i := range ur {
		req, err := http.NewRequest(http.MethodGet, ur[i], nil)
		if err != nil {
			return err
		}

		req.Header.Add("referer", "https://readmanganato.com/")

		for retries := 7; retries > 0; retries-- {
			res, err := client.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.StatusCode >= 400 && res.StatusCode < 500 {
				return errors.New(res.Status)
			}

			if res.StatusCode == 200 {
				f, err := os.Create(dirName + "/Documents/" + name + "/" + strconv.Itoa(chap) + "/" + strconv.Itoa(i) + ".jpg")
				if err != nil {
					return err
				}
				defer f.Close()

				_, err = io.Copy(f, res.Body)
				if err != nil {
					return err
				}
			}
		}
	}

	return errors.New("too much retries")
}

func (m Manga) DownloadMultiFromPage(start, end int) {
	length := len(m.Chapter)
	x := start
	for i := length - start; i >= length-end; i-- {
		s := GetMangaPicLink(m.Chapter[i].Url)
		DownloadFromPage(s, x, m.Title)
		x += 1
	}
}
