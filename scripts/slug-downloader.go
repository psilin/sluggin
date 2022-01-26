package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	URL              string = "https://support.allizom.org/api/1/kb/"
	MaxChannelLength int    = 16
)

type NameResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"prev"`
	Slugs    []struct {
		Title string `json:"title"`
		URL   string `json:"slug"`
	} `json:"results"`
}

func getSlugNames(out chan string, verbose bool, num int) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result NameResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}
	actualNum := result.Count
	if num > actualNum {
		fmt.Printf("Total number of slugs: %v is less than what we want: %v, will download what we can.", actualNum, num)
		num = actualNum
	}

	cnt := 0
	for _, s := range result.Slugs {
		if cnt == num {
			break
		}
		out <- s.URL
		if verbose {
			fmt.Printf("%v added slug %v to task queue\n", cnt, s.URL)
		}
		cnt += 1
	}
	resp.Body.Close()

	for cnt < num {
		resp, err := http.Get(result.Next)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var result NameResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalln(err)
		}

		for _, s := range result.Slugs {
			if cnt == num {
				break
			}
			out <- s.URL
			if verbose {
				fmt.Printf("%v added slug %v to task queue\n", cnt, s.URL)
			}
			cnt += 1
		}
		resp.Body.Close()
	}
	close(out)
}

func DownloadSlugs(verbose bool, num int, path string) {
	in := make(chan string, MaxChannelLength)
	go getSlugNames(in, verbose, num)
	cnt := 0
	for sl := range in {
		fmt.Printf("%v %v\n", cnt, sl)
		cnt += 1
	}
}
