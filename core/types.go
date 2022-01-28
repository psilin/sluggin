package core

type Slug struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Slug     string   `json:"slug"`
	Url      string   `json:"url"`
	Locale   string   `json:"locale"`
	Products []string `json:"products"`
	Topics   []string `json:"topics"`
	Summary  string   `json:"summary"`
}

type SlugNames struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"prev"`
	Slugs    []struct {
		Title string `json:"title"`
		URL   string `json:"slug"`
	} `json:"results"`
}
