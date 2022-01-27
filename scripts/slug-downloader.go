package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/psilin/sluggin/db"
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

type SlugResponse struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Slug     string   `json:"slug"`
	Url      string   `json:"url"`
	Locale   string   `json:"locale"`
	Products []string `json:"products"`
	Topics   []string `json:"topics"`
	Summary  string   `json:"summary"`
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

func processSlug(verbose bool, idx int, in chan string, out chan [2]int, dbase *sqlx.DB) {
	errors := 0
	successes := 0
	for s := range in {
		if verbose {
			fmt.Printf("%v started processing %v\n", idx, s)
		}
		errors += 1
		resp, err := http.Get(URL + s)
		if err != nil {
			fmt.Printf("Error getting: %v\n", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading: %v\n", err)
			continue
		}

		var result SlugResponse
		if err := json.Unmarshal(body, &result); err != nil {
			resp.Body.Close()
			fmt.Printf("Error unmarshaling: %v\n", err)
			continue
		}

		// FIX ME insert into DB here

		if verbose {
			fmt.Printf("%v finished processing %v\n", idx, s)
		}
		errors -= 1
		successes += 1
	}
	out <- [2]int{successes, errors}
}

func DownloadSlugs(verbose bool, num int) {
	// FIX ME extract connection string from here
	dbase, err := db.InitDb("user=postgres password=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.InitialCleanup(dbase)
	if err != nil {
		log.Fatalln(err)
	}

	in := make(chan string, MaxChannelLength)
	out := make(chan [2]int, runtime.NumCPU())
	go getSlugNames(in, verbose, num)

	if verbose {
		fmt.Printf("Detected %v CPUs.\n", runtime.NumCPU())
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go processSlug(verbose, i, in, out, dbase)
	}

	errors := 0
	successes := 0
	for i := 0; i < runtime.NumCPU(); i++ {
		// use this instead of wait group
		res := <-out
		successes += res[0]
		errors += res[1]
	}
	fmt.Printf("Processed %v slugs of %v, success: %v, error: %v\n", successes+errors, num, successes, errors)

	err = db.CloseDb(dbase)
	if err != nil {
		log.Fatalln(err)
	}
}
