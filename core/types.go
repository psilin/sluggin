package core

type Slug struct {
	Id       int      `json:"id" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Slug     string   `json:"slug" binding:"required"`
	Url      string   `json:"url" binding:"required"`
	Locale   string   `json:"locale" binding:"required"`
	Products []string `json:"products" binding:"required"`
	Topics   []string `json:"topics" binding:"required"`
	Summary  string   `json:"summary" binding:"required"`
}

type OutSlug struct {
	DbId int  `json:"dbid"`
	Slg  Slug `json:"slug"`
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
